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
	/*
		for each obj in the perspective:
			if obj in p1 but not in p2: add to p2
			elif obj in p2 but not in p1: add to p1
			else:
				for each version of the object in p1 and p2:
					if version and hash match: continue
					elif: conflict, error out
				if done with p1 but not p2:
					add all remaining versions to p1
				if done with p2 but not p1:
					add all remaining versions to p2
	*/
}
