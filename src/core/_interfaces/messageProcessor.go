package _interfaces

type MessageProcessor interface {
	Process(message Message) (string, error)
}
