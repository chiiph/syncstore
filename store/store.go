package store

import (
	"fmt"
	"syncstore/perspective"

	"github.com/pkg/errors"
)

const PerspectiveKey = "Perspective"

var ErrObjectNotFound = errors.New("object not found")

// ObjectStore is a simple interface to an object store, could be s3, minio, or simple files
type ObjectStore interface {
	Put(key string, data []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
	List() chan string
}

// Store is a perspective enriched / versioned object store
type Store interface {
	ObjectStore
	PutWithVersion(key string, version int, data []byte) error
	Perspective() perspective.Perspective
}

func New(os ObjectStore, prefix string) Store {
	p := perspective.New()

	rawP, err := os.Get(key(PerspectiveKey, prefix, 0))
	if err != nil {
		if err != ErrObjectNotFound {
			panic(err)
		}
	}
	if rawP != nil {
		err := p.Unmarshal(rawP)
		if err != nil {
			panic(err)
		}
	}

	return &storeImpl{
		p:      p,
		os:     os,
		prefix: prefix,
	}
}

type storeImpl struct {
	p  perspective.Perspective
	os ObjectStore

	prefix string
}

func key(k string, prefix string, version int) string {
	key := ""
	if prefix != "" {
		key = prefix + "/"
	}
	key += k
	if version > 0 {
		key = fmt.Sprintf("%s/%d", key, version)
	}
	return key
}

func (s *storeImpl) key(k string, version int) string {
	return key(k, s.prefix, version)
}

func (s *storeImpl) PutWithVersion(key string, version int, data []byte) error {
	err := s.os.Put(s.key(key, version), data)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.p.Update(key, version)
	if err != nil {
		return errors.WithStack(err)
	}

	// TODO: store perspective in os or rollback version if it fails

	return nil
}

func (s *storeImpl) Put(key string, data []byte) error {
	v, err := s.p.Version(key)
	if err != nil {
		if err != perspective.ErrNoSuchObject {
			return err
		}
	}

	v += 1

	return s.PutWithVersion(key, v, data)
}

func (s *storeImpl) Get(key string) ([]byte, error) {
	v, err := s.p.Version(key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	data, err := s.os.Get(s.key(key, v))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func (s *storeImpl) Del(key string) error {
	panic("Not implemented")
}

func (s *storeImpl) List() chan string {
	panic("Not implemented")
	return nil
}

func (s *storeImpl) Perspective() perspective.Perspective {
	return s.p
}
