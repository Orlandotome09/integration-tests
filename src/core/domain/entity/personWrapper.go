package entity

import "sync"

type PersonWrapper struct {
	Person Person
	Mutex  sync.Mutex
}
