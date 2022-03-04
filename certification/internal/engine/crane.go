package engine

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
)

type craneEngine struct{}

func NewCraneEngine() *craneEngine {
	return &craneEngine{}
}

func (c *craneEngine) ListTags(imageURI string) ([]string, error) {
	// prepare crane runtime options, if necessary
	options := make([]crane.Option, 0, 1)
	options = append(options, crane.WithAuthFromKeychain(authn.DefaultKeychain))

	return crane.ListTags(imageURI, options...)
}
