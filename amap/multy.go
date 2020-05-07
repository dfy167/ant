package amap

import "sync"

// MRun MRun
type MRun struct {
	c  chan struct{}
	wg *sync.WaitGroup
}

// NewMRun NewMRun
func NewMRun(maxSize int) *MRun {
	return &MRun{
		c:  make(chan struct{}, maxSize),
		wg: new(sync.WaitGroup),
	}
}

// Add Add
func (s *MRun) Add(delta int) {
	s.wg.Add(delta)
	for i := 0; i < delta; i++ {
		s.c <- struct{}{}
	}
}

// Done Done
func (s *MRun) Done() {
	<-s.c
	s.wg.Done()
}

// Wait Wait
func (s *MRun) Wait() {
	s.wg.Wait()
}
