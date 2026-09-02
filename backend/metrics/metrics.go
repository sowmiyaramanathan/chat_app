package metrics

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type durationMetric struct {
	mu      sync.Mutex
	samples []time.Duration
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(body)
}

func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)
		Observe(fmt.Sprintf("http.%s.%s", r.Method, r.URL.Path), time.Since(start))
	})
}

var (
	durations       sync.Map
	connections     atomic.Int64
	peakConnections atomic.Int64
	messagesRouted  atomic.Int64
	bufferFull      atomic.Int64
	poolStatsMu     sync.RWMutex
	lastPoolStats   sql.DBStats
	poolPeakOpen    int
	poolPeakInUse   int
	lastRuntime     RuntimeStats
)

type RuntimeStats struct {
	Goroutines     int
	PeakGoroutines int
	HeapAlloc      uint64
	HeapInUse      uint64
	NumGC          uint32
	CPUPercent     float64
	PeakCPUPercent float64
}

func Observe(name string, d time.Duration) {
	value, _ := durations.LoadOrStore(name, &durationMetric{})
	m := value.(*durationMetric)
	m.mu.Lock()
	m.samples = append(m.samples, d)
	if len(m.samples) > 100000 {
		m.samples = append([]time.Duration(nil), m.samples[len(m.samples)-50000:]...)
	}
	m.mu.Unlock()
}

func IncMessagesRouted() { messagesRouted.Add(1) }
func IncBufferFull()     { bufferFull.Add(1) }

func ConnectionOpened() {
	current := connections.Add(1)
	for {
		peak := peakConnections.Load()
		if current <= peak || peakConnections.CompareAndSwap(peak, current) {
			return
		}
	}
}

func ConnectionClosed() { connections.Add(-1) }

func StartRuntimeSampler(db *sql.DB) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		updateRuntime(RuntimeStats{Goroutines: runtime.NumGoroutine(), HeapAlloc: mem.HeapAlloc, HeapInUse: mem.HeapInuse, NumGC: mem.NumGC})
		if db != nil {
			poolStatsMu.Lock()
			lastPoolStats = db.Stats()
			updatePoolPeaksLocked(lastPoolStats)
			poolStatsMu.Unlock()
		}
		lastCPU := processCPUTime()
		lastWall := time.Now()

		for range ticker.C {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			now := time.Now()
			cpuTime := processCPUTime()
			wall := now.Sub(lastWall).Seconds()
			cpu := 0.0
			if wall > 0 {
				cpu = (cpuTime - lastCPU).Seconds() / wall * 100
			}
			lastCPU, lastWall = cpuTime, now
			if cpu < 0 {
				cpu = 0
			}

			updateRuntime(RuntimeStats{Goroutines: runtime.NumGoroutine(), HeapAlloc: mem.HeapAlloc, HeapInUse: mem.HeapInuse, NumGC: mem.NumGC, CPUPercent: cpu})
			if db != nil {
				poolStatsMu.Lock()
				lastPoolStats = db.Stats()
				updatePoolPeaksLocked(lastPoolStats)
				poolStatsMu.Unlock()
			}
		}
	}()
}

func updatePoolPeaksLocked(stats sql.DBStats) {
	if stats.OpenConnections > poolPeakOpen {
		poolPeakOpen = stats.OpenConnections
	}
	if stats.InUse > poolPeakInUse {
		poolPeakInUse = stats.InUse
	}
}

func updateRuntime(stats RuntimeStats) {
	poolStatsMu.Lock()
	if stats.Goroutines < lastRuntime.PeakGoroutines {
		stats.PeakGoroutines = lastRuntime.PeakGoroutines
	}
	if stats.CPUPercent < lastRuntime.PeakCPUPercent {
		stats.PeakCPUPercent = lastRuntime.PeakCPUPercent
	}
	if stats.CPUPercent > stats.PeakCPUPercent {
		stats.PeakCPUPercent = stats.CPUPercent
	}
	lastRuntime = stats
	poolStatsMu.Unlock()
}

func Snapshot() string {
	poolStatsMu.RLock()
	pool, runtimeStats := lastPoolStats, lastRuntime
	poolStatsMu.RUnlock()

	return fmt.Sprintf("DATABASE/POOL open=%d in_use=%d idle=%d peak_open=%d peak_in_use=%d wait_count=%d wait_duration=%s\nAPPLICATION goroutines=%d peak_goroutines=%d heap_alloc=%s heap_in_use=%s gc=%d cpu=%.1f%% peak_cpu=%.1f%%\nWEBSOCKET connections=%d peak_connections=%d messages_routed=%d buffer_full=%d\nTIMINGS\n%s",
		pool.OpenConnections, pool.InUse, pool.Idle, poolPeakOpen, poolPeakInUse, pool.WaitCount, pool.WaitDuration,
		runtimeStats.Goroutines, runtimeStats.PeakGoroutines, formatBytes(runtimeStats.HeapAlloc), formatBytes(runtimeStats.HeapInUse), runtimeStats.NumGC, runtimeStats.CPUPercent, runtimeStats.PeakCPUPercent,
		connections.Load(), peakConnections.Load(), messagesRouted.Load(), bufferFull.Load(), timingSummary())
}

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprint(w, Snapshot())
}

func timingSummary() string {
	var names []string
	durations.Range(func(key, _ any) bool { names = append(names, key.(string)); return true })
	sort.Strings(names)
	result := ""
	for _, name := range names {
		value, _ := durations.Load(name)
		m := value.(*durationMetric)
		m.mu.Lock()
		samples := append([]time.Duration(nil), m.samples...)
		m.mu.Unlock()
		if len(samples) == 0 {
			continue
		}
		sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
		result += fmt.Sprintf("%s count=%d p50=%s p95=%s p99=%s\n", name, len(samples), percentile(samples, .50), percentile(samples, .95), percentile(samples, .99))
	}
	return result
}

func percentile(values []time.Duration, p float64) time.Duration {
	return values[int(float64(len(values)-1)*p)]
}

func formatBytes(value uint64) string {
	if value < 1024 {
		return fmt.Sprintf("%dB", value)
	}
	return fmt.Sprintf("%.1fMiB", float64(value)/(1024*1024))
}

func processCPUTime() time.Duration {
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		return 0
	}
	return time.Duration(usage.Utime.Sec)*time.Second + time.Duration(usage.Utime.Usec)*time.Microsecond + time.Duration(usage.Stime.Sec)*time.Second + time.Duration(usage.Stime.Usec)*time.Microsecond
}
