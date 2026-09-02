package utils

import (
	"backend/entities/packet"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func EncodeCursor(time time.Time, id uint) string {
	payload := fmt.Sprintf("%d_%d", time.UnixNano(), id)
	return base64.URLEncoding.EncodeToString([]byte(payload))
}

func DecodeCursor(encoded string) (cursor *packet.Cursor, err error) {
	if encoded == "" {
		return
	}

	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return
	}

	parts := strings.Split(string(decoded), "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cursor format")
	}

	nanos, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp in cursor: %w", err)
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid id in cursor: %w", err)
	}

	return &packet.Cursor{
		CreatedAt: time.Unix(0, nanos),
		ID:        id,
	}, nil
}
