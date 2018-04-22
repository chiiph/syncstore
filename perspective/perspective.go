package perspective

import (
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
)

// Item is an item of the perspective, this is used for iterating and bulk updates of items
type Item struct {
	Obj     string
	Version int
}

var ErrNoSuchObject = errors.New("object does not exist")
var ErrInvalidVersion = errors.New("version should be greater than 0")

// Perspective is a simple way to keep track of versions of objects in a store
// It is used by the synchronization algorithm to produce a plan for keeping
// stores up to date
type Perspective interface {
	// Update adds or overwrites the last version for a specified object
	Update(obj string, version int) error
	// Version returns the version for the object or an error if the object is unknown
	Version(obj string) (int, error)
	// Marshal serializes the perspective to be stored
	Marshal() ([]byte, error)
	// Unmarshal populates the perspective with the data from the serialized version provided
	Unmarshal(b []byte) error
	// Iterate gives an easy way to go through the whole perspective
	Iterate() chan Item
}

func New() Perspective {
	return &perspectiveImpl{
		items: make(map[string]int),
	}
}

type perspectiveImpl struct {
	l     sync.Mutex
	items map[string]int
}

func (p *perspectiveImpl) Update(obj string, version int) error {
	p.l.Lock()
	defer p.l.Unlock()

	if version <= 0 {
		return ErrInvalidVersion
	}

	// TODO: Prevent version rollbacks

	p.items[obj] = version
	return nil
}

func (p *perspectiveImpl) Version(obj string) (int, error) {
	p.l.Lock()
	defer p.l.Unlock()

	v, ok := p.items[obj]
	if !ok {
		return 0, ErrNoSuchObject
	}

	return v, nil
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

func (p *perspectiveImpl) Iterate() chan Item {
	items := make(chan Item)

	go func() {
		p.l.Lock()
		defer p.l.Unlock()

		for k, v := range p.items {
			p.l.Unlock()
			items <- Item{k, v}
			p.l.Lock()
		}
	}()

	return items
}
