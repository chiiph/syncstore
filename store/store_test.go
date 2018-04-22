package store

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func newS() Store {
	os := NewFakeOS()
	return New(os, "acc1")
}

func TestBasic(t *testing.T) {
	s := newS()
	require.NoError(t, s.Put("k1", []byte{1, 2, 3}))
	d, e := s.Get("k1")
	require.NoError(t, e)
	require.Equal(t, []byte{1, 2, 3}, d)
}

func checkVersion(t *testing.T, s Store, key string, expected int) {
	p := s.Perspective()
	v, e := p.Version(key)
	require.NoError(t, e)
	require.Equal(t, expected, v)
}

func TestPerspective(t *testing.T) {
	s := newS()
	require.NoError(t, s.Put("k1", []byte{1, 2, 3}))
	checkVersion(t, s, "k1", 1)

	require.NoError(t, s.Put("k1", []byte{4}))
	checkVersion(t, s, "k1", 2)

	require.NoError(t, s.Put("k2", []byte{1, 2, 3}))
	checkVersion(t, s, "k2", 1)

	// TODO: test perspective unmarshaling from an existent store
}

type fakeOS struct {
	l     sync.Mutex
	items map[string][]byte
}

func NewFakeOS() *fakeOS {
	return &fakeOS{
		items: make(map[string][]byte),
	}
}

func (f *fakeOS) Put(key string, data []byte) error {
	f.l.Lock()
	defer f.l.Unlock()
	f.items[key] = data
	return nil
}

func (f *fakeOS) Get(key string) ([]byte, error) {
	f.l.Lock()
	defer f.l.Unlock()
	data, ok := f.items[key]
	if !ok {
		return nil, ErrObjectNotFound
	}
	return data, nil
}

func (f *fakeOS) Del(key string) error {
	f.l.Lock()
	defer f.l.Unlock()
	_, ok := f.items[key]
	if !ok {
		return ErrObjectNotFound
	}
	delete(f.items, key)
	return nil
}

func (f *fakeOS) List() chan string {
	f.l.Lock()
	defer f.l.Unlock()
	return nil
}
