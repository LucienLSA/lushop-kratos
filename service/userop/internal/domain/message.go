package domain

type Message struct {
	ID          int32
	UserID      int32
	MessageType int32
	Subject     string
	Message     string
	File        string
}

type MessageListResponse struct {
	Total int64
	List  []*Message
}
