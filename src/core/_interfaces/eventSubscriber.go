package _interfaces

type EventSubscriber interface {
	Process() error
}
