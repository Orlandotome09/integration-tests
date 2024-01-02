package _interfaces

type GrpcManager interface {
	Listen(port int, cncInstance EventProcessor) error
}
