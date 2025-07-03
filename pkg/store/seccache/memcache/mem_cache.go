package memcache

import (
	"errors"
	"sync"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
)

// TODO UMU Need to make this concurrent-safe (lock...)
type MemSecureCache struct {
	db map[string]*did.Doc
	sync.RWMutex
}

func New() (*MemSecureCache, error) {
	return &MemSecureCache{db: make(map[string]*did.Doc)}, nil
}

func (c *MemSecureCache) StoreDidDoc(document *did.Doc) error {
	if document == nil {
		return errors.New("value cannot be nil")
	}
	if document.ID == "" {
		return errors.New("didId cannot be empty")
	}
	c.Lock()
	defer c.Unlock()
	c.db[document.ID] = document
	return nil
}

func (c *MemSecureCache) RetrieveDidDoc(didID string) (bool, *did.Doc, error) {
	if didID == "" {
		return false, nil, errors.New("didId cannot be empty")
	}
	c.RLock()
	defer c.RUnlock()
	entry, ok := c.db[didID]

	if !ok {
		return false, nil, nil
	}

	return true, entry, nil
}

func (c *MemSecureCache) RemoveAll() error {
	for k := range c.db {
		delete(c.db, k)
	}

	return nil
}
