package cr

import (
	"context"

	"github.com/frantjc/kontrol/pkg"
)

type ContainerRegistry interface {
	Pull(context.Context, string) (pkg.Image, error)
	Push(context.Context, string, pkg.Image) error
}
