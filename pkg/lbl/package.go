package lbl

import (
	"context"

	"github.com/frantjc/kontrol/pkg"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"golang.org/x/exp/maps"
)

const (
	LabelPackaging = "cc.frantj.kontrol.packaging"
)

func (p *Packager) Package(_ context.Context, image pkg.Image, packaging *pkg.Packaging) (pkg.Image, error) {
	cfgf, err := image.ConfigFile()
	if err != nil {
		return nil, err
	}

	cfg := cfgf.Config.DeepCopy()
	if cfg.Labels == nil {
		cfg.Labels = map[string]string{}
	}

	value, err := p.Encode(packaging)
	if err != nil {
		return nil, err
	}

	maps.Copy(cfg.Labels, map[string]string{
		LabelPackaging: value,
	})

	return mutate.Config(image, *cfg)
}
