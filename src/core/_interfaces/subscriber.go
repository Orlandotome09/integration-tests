package _interfaces

type Subscriber interface {
	Listen(processor Processor) error
}
