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

type MessageEvent struct {
	EventID        string    `json:"event_id"`
	MessageID      string    `json:"message_id"`
	SenderID       string    `json:"sender_id"`
	RecipientID    string    `json:"recipient_id"`
	Payload        []byte    `json:"payload"`
	OriginInstance string    `json:"origin_instance"`
	CreatedAt      time.Time `json:"created_at"`
}
