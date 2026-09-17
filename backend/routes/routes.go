package routes

import (
	"backend/auth"
	"backend/controllers"
	"backend/metrics"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

func InitializeRoutes(c controllers.Controller) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(slogRequestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   frontendOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Content-type", "Authorization"},
		AllowCredentials: true,
	}))

	r.HandleFunc("/ws/{userID}", c.HandleConnection)

	r.Get("/metrics", metrics.Handler)
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", c.RegisterUser)
		r.Post("/login", c.LoginUser)
		r.Post("/auth/refresh", c.RefreshToken)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(auth.TokenAuth))
			r.Use(jwtauth.Authenticator(auth.TokenAuth))
			r.Get("/profile", c.Profile)
			r.Get("/friends", c.GetMyFriends)
			r.Get("/users", c.GetNonFriends)
		})
	})

	r.Route("/message", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(metrics.HTTPMiddleware)
			r.Use(jwtauth.Verifier(auth.TokenAuth))
			r.Use(jwtauth.Authenticator(auth.TokenAuth))
			r.Post("/create", c.CreateMessage)
			r.Get("/view", c.GetMessages)
		})
	})

	r.Route("/friends", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(auth.TokenAuth))
			r.Use(jwtauth.Authenticator(auth.TokenAuth))
			r.Get("/isFriend", c.IsFriend)
			r.Post("/sendRequest", c.SendFriendRequest)
			r.Get("/getFriendRequests", c.GetFriendRequests)
			r.Put("/acceptRequest", c.AcceptFriendRequest)
			r.Put("/rejectRequest", c.RejectFriendRequest)
			r.Get("/isRequestSent", c.IsFriendRequestSent)
			r.Get("/isRequestReceived", c.IsRequestReceived)
		})
	})

	return r
}

func frontendOrigins() []string {
	configured := os.Getenv("FRONTEND_ORIGINS")
	if configured == "" {
		configured = os.Getenv("ALLOWED_ORIGIN")
	}
	if configured == "" {
		if os.Getenv("ENV") == "PROD" {
			slog.Warn("FRONTEND_ORIGINS/ALLOWED_ORIGIN must be set in production")
			return []string{}
		}

		return []string{"http://localhost:3000", "http://localhost:3001"}
	}

	origins := strings.Split(configured, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}

func ValidateProductionConfig() error {
	if os.Getenv("ENV") != "PROD" {
		return nil
	}
	origins := frontendOrigins()
	if len(origins) == 0 {
		return fmt.Errorf("FRONTEND_ORIGINS or ALLOWED_ORIGIN must be set in production")
	}
	for _, origin := range origins {
		if strings.Contains(origin, "localhost") {
			return fmt.Errorf("localhost must not be configured as a production frontend origin")
		}
	}
	return nil
}

// slogRequestLogger records request metadata without exposing query parameters,
// which may contain credentials such as the WebSocket token.
func slogRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}
