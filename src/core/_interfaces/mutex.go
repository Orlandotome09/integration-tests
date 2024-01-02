package _interfaces

type Mutex interface {
	Lock(id string) (success bool)
	Release(id string)
}
