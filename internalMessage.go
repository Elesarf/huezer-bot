package main

// InternalMessage внутреннее представление сообщения
type InternalMessage struct {
	chatID      int64
	messageID   int
	userName    string
	messageText string
}
