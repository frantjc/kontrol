package cr

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/frantjc/kontrol/pkg"
)

var Default = func() ContainerRegistry {
	crs := []ContainerRegistry{}

	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err == nil {
		if _, err = c.Ping(context.Background()); err == nil {
			crs = append(crs, &Daemon{c})
		}
	}

	return Coalesce(append(crs, new(Remote)))
}()

func Pull(ctx context.Context, ref string) (pkg.Image, error) {
	return Default.Pull(ctx, ref)
}

func Push(ctx context.Context, ref string, image pkg.Image) error {
	return Default.Push(ctx, ref, image)
}
