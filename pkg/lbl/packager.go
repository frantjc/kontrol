package lbl

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
)

type Packager struct {
	Encoder
}

type Encoder interface {
	Encode(any) (string, error)
	Decode(string, any) error
}

func (e *Packager) Encode(v any) (string, error) {
	if e.Encoder != nil {
		return e.Encoder.Encode(v)
	}

	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(buf.Bytes()), nil
}

func (e *Packager) Decode(s string, v any) error {
	if e.Encoder != nil {
		return e.Encoder.Decode(s, v)
	}

	return json.NewDecoder(
		base64.NewDecoder(
			base64.RawStdEncoding,
			strings.NewReader(s),
		),
	).Decode(v)
}
