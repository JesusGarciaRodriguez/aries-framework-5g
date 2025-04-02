package seccache

import "github.com/hyperledger/aries-framework-go/pkg/doc/did"

// XXX UMU: This interface and implementation will be used for integrating secure cache as needed (in 5g development) for simplicity's sake
// The same functionality could be achieved with carefull usage and implementation of Storage interfaces already in ARIES, e.g. DIDStore, mem and creating a new one for the secure cache

type SecureCache interface {
	//Stores Did document (using its ID as index)
	StoreDidDoc(document *did.Doc) error
	//Retrieves Did document for a specific DID. It assumes (as other ARIES interfaces used) that the didID identifies the document, and not a specific key inside
	RetrieveDidDoc(didID string) (bool, *did.Doc, error)
}

type NoOpSecureCache struct{}

func NewNoOpCache() (*NoOpSecureCache, error) {
	return &NoOpSecureCache{}, nil
}

func (c *NoOpSecureCache) StoreDidDoc(document *did.Doc) error {
	return nil
}

func (c *NoOpSecureCache) RetrieveDidDoc(didID string) (bool, *did.Doc, error) {
	return false, nil, nil
}
