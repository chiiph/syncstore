package ossync

import (
	"syncstore/perspective"
	"syncstore/store"
)

// OSSync is an object store synchronization helper
// Given 2 Stores, it uses their perspective to sync one with the other
type OSSync interface {
	Sync() error
}

func New(from, to store.Store) OSSync {
	return &ssyncImpl{}
}

type ssyncImpl struct{}

func (s *ssyncImpl) Sync() error {
	return nil
}

// diff produces a set of steps to get p2 up to date with p1
func (s *ssyncImpl) diff(p1, p2 perspective.Perspective) {

}
