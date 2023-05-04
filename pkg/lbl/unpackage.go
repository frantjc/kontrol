package lbl

import (
	"context"

	"github.com/frantjc/kontrol/pkg"
)

func (p *Packager) Unpackage(_ context.Context, image pkg.Image) (*pkg.Packaging, error) {
	cfgf, err := image.ConfigFile()
	if err != nil {
		return nil, err
	}

	var (
		cfg       = cfgf.Config.DeepCopy()
		packaging = &pkg.Packaging{}
	)

	value, ok := cfg.Labels[LabelPackaging]
	if ok {
		if err = p.Decode(value, packaging); err != nil {
			return nil, err
		}
	}

	return packaging, nil
}
