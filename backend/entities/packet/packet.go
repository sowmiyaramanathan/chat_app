package packet

import "time"

type Users struct {
	ID       uint64
	Username string
}

type Messages struct {
	ID         uint
	FromUserID uint64
	ToUserID   uint64
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
	FromUserID uint64
	Username   string
}

type Cursor struct {
	CreatedAt time.Time
	ID        int64
}
