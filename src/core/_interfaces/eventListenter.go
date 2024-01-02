package _interfaces

type EventListener interface {
	Register(subscriber QueueSubscriber, processor MessageProcessor)
	Listen()
}
