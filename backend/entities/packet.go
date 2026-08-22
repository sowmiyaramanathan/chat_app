package entities

type Users struct {
	ID       uint64
	Username string
}
type Messages struct {
	FromUserID uint64
	ToUserID   uint64
	Message    string
}
type Requests struct {
	FromUserID uint64
	Username   string
}
