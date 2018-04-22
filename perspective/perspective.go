package perspective

import (
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
)

// Item is an item of the perspective, this is used for iterating and bulk updates of items
type Item struct {
	Obj     string `json:"key"`
	Version int    `json:"version"`
	Hash    []byte `json:"hash"`
}

var ErrNoSuchObject = errors.New("object does not exist")
var ErrInvalidVersion = errors.New("version should be greater than 0")

// Perspective is a simple way to keep track of versions of objects in a store
// It is used by the synchronization algorithm to produce a plan for keeping
// stores up to date
type Perspective interface {
	// Update adds or overwrites the last version for a specified object
	Update(obj string, version int, hash []byte) error
	// Version returns the version for the object or an error if the object is unknown
	Version(obj string) (int, error)
	// Marshal serializes the perspective to be stored
	Marshal() ([]byte, error)
	// Unmarshal populates the perspective with the data from the serialized version provided
	Unmarshal(b []byte) error
	// Iterate gives an easy way to go through the whole perspective
	Iterate() chan *Item
}

func New() Perspective {
	return &perspectiveImpl{
		items: make(map[string][]*Item),
	}
}

type perspectiveImpl struct {
	l     sync.Mutex
	items map[string][]*Item
}

func (p *perspectiveImpl) Update(obj string, version int, hash []byte) error {
	p.l.Lock()
	defer p.l.Unlock()

	if version <= 0 {
		return ErrInvalidVersion
	}

	// TODO: Prevent version rollbacks

	br := p.items[obj]
	it := &Item{
		Obj:     obj,
		Version: version,
		Hash:    hash,
	}
	br = append(br, it)
	p.items[obj] = br

	return nil
}

func (p *perspectiveImpl) Version(obj string) (int, error) {
	p.l.Lock()
	defer p.l.Unlock()

	i, ok := p.items[obj]
	if !ok {
		return 0, ErrNoSuchObject
	}

	return i[len(i)-1].Version, nil
}

func (p *perspectiveImpl) Marshal() ([]byte, error) {
	p.l.Lock()
	defer p.l.Unlock()

	b, err := json.Marshal(p.items)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return b, nil
}

func (p *perspectiveImpl) Unmarshal(b []byte) error {
	p.l.Lock()
	defer p.l.Unlock()

	err := json.Unmarshal(b, p.items)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (p *perspectiveImpl) Iterate() chan *Item {
	items := make(chan *Item)

	go func() {
		p.l.Lock()
		defer p.l.Unlock()

		for _, v := range p.items {
			p.l.Unlock()
			items <- v[len(v)-1]
			p.l.Lock()
		}
	}()

	return items
}

func (p *perspectiveImpl) IterateObj(obj string) chan *Item {
	items := make(chan *Item)

	go func() {
		p.l.Lock()
		defer p.l.Unlock()

		for _, v := range p.items[obj] {
			p.l.Unlock()
			items <- v
			p.l.Lock()
		}
	}()

	return items
}
