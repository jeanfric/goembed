// Package base64embedder implements an asset embedder that encodes
// assets as base64 strings.
package base64embedder

import (
	"encoding/base64"
	"io"
	"io/ioutil"

	"github.com/jeanfric/embed"
)

var (
	imports = [...]string{"encoding/base64"}
)

const (
	decode = `func(s string) (string, error) {
		b, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}`
)

func encode(contents io.Reader) (string, error) {
	b, err := ioutil.ReadAll(contents)
	if err != nil {
		return "", err
	}
	return "`" + base64.StdEncoding.EncodeToString(b) + "`", nil
}

// NewSequential creates a new sequential base64embedder asset embedder.
func NewSequential() embed.AssetEmbedder {
	return embed.NewSequentialEmbedder(encode, decode, imports[:])
}

// NewConcurrent creates a new concurrent base64embedder asset embedder.
func NewConcurrent() embed.AssetEmbedder {
	return embed.NewConcurrentEmbedder(encode, decode, imports[:])
}
