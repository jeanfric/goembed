// Package hexembedder implements an asset embedder that encodes
// assets as hexadecimal strings.
package hexembedder

import (
	"encoding/hex"
	"io"
	"io/ioutil"

	"github.com/jeanfric/embed"
)

var (
	imports = [...]string{"encoding/hex"}
)

const (
	decode = `func(s string) (string, error) {
		b, err := hex.DecodeString(s)
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
	return "`" + hex.EncodeToString(b) + "`", nil
}

// NewSequential creates a new sequential hexembedder asset embedder.
func NewSequential() embed.AssetEmbedder {
	return embed.NewSequentialEmbedder(encode, decode, imports[:])
}

// NewConcurrent creates a new concurrent hexembedder asset embedder.
func NewConcurrent() embed.AssetEmbedder {
	return embed.NewConcurrentEmbedder(encode, decode, imports[:])
}
