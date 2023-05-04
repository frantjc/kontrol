package cr

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/frantjc/kontrol/pkg"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
)

type DaemonClient interface {
	daemon.Client
	ImageInspectWithRaw(context.Context, string) (types.ImageInspect, []byte, error)
	ImagePull(context.Context, string, types.ImagePullOptions) (io.ReadCloser, error)
}

type Daemon struct {
	DaemonClient
}

func (cr *Daemon) Pull(ctx context.Context, ref string) (pkg.Image, error) {
	reference, err := name.NewTag(ref)
	if err != nil {
		return nil, err
	}

	if _, _, err = cr.ImageInspectWithRaw(ctx, ref); err != nil {
		ipr, err := cr.ImagePull(ctx, ref, types.ImagePullOptions{})
		if err != nil {
			return nil, err
		}

		if _, err = io.Copy(io.Discard, ipr); err != nil {
			return nil, err
		}

		if err = ipr.Close(); err != nil {
			return nil, err
		}
	}

	return daemon.Image(reference, daemon.WithContext(ctx), daemon.WithClient(cr))
}

func (cr *Daemon) Push(ctx context.Context, ref string, image pkg.Image) error {
	tag, err := name.NewTag(ref)
	if err != nil {
		return err
	}

	_, err = daemon.Write(tag, image, daemon.WithContext(ctx), daemon.WithClient(cr))
	return err
}
