package websocket

type MessageType string

const (
	MessageTypeNotification MessageType = "notification"
	MessageTypeMessage      MessageType = "message"
	MessageTypeChat         MessageType = "chat"
)

func (m MessageType) Value() string {
	return string(m)
}
