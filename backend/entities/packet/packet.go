package packet

import "time"

type Users struct {
	ID       string
	UserName string
}

type Messages struct {
	ID         uint
	FromUserID string
	ToUserID   string
	Message    string
	CreatedAt  time.Time
}

type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}

type AllMessagesRes struct {
	Messages []*Messages
	PageInfo PageInfo
}

type Requests struct {
	FromUserID string
	UserName   string
}

type Cursor struct {
	CreatedAt time.Time
	ID        int64
}

type ErrorResponse struct {
	Error string `json:"error"`
}
