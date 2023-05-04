package pkg

import (
	"context"
)

type Packager interface {
	Package(context.Context, Image, *Packaging) (Image, error)
	Unpackage(context.Context, Image) (*Packaging, error)
}
