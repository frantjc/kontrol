package cr

import (
	"context"

	"github.com/frantjc/kontrol/pkg"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type Remote struct{}

func (cr *Remote) Pull(ctx context.Context, ref string) (pkg.Image, error) {
	reference, err := name.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	return remote.Image(reference, remote.WithContext(ctx), remote.WithAuthFromKeychain(authn.DefaultKeychain))
}

func (cr *Remote) Push(ctx context.Context, ref string, image pkg.Image) error {
	tag, err := name.NewTag(ref)
	if err != nil {
		return err
	}

	return remote.Write(tag, image, remote.WithContext(ctx))
}
