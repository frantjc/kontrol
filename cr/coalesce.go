package cr

import (
	"context"
	"errors"

	"github.com/frantjc/kontrol/pkg"
)

type Coalesce []ContainerRegistry

func (c Coalesce) Pull(ctx context.Context, ref string) (pkg.Image, error) {
	errs := []error{}

	for _, cr := range c {
		image, err := cr.Pull(ctx, ref)
		if err == nil {
			return image, nil
		}

		errs = append(errs, err)
	}

	return nil, errors.Join(errs...)
}

func (c Coalesce) Push(ctx context.Context, ref string, image pkg.Image) error {
	errs := []error{}

	for _, cr := range c {
		err := cr.Push(ctx, ref, image)
		if err == nil {
			return nil
		}

		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
