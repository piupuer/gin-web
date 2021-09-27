package ch

import (
	"sync"
)

type Ch struct {
	C      chan interface{}
	closed bool
	lock   sync.Mutex
}

func NewCh() *Ch {
	return &Ch{C: make(chan interface{})}
}

func (s *Ch) SafeClose() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.closed {
		close(s.C)
		s.closed = true
	}
}

func (s Ch) IsClosed() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.closed
}

func (s Ch) SafeSend(data interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.closed {
		s.C <- data
	}
}
