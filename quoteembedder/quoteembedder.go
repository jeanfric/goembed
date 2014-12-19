// Package quoteembedder implements an asset embedder that encodes
// assets as quoted strings.
package quoteembedder

import (
	"io"
	"io/ioutil"
	"strconv"

	"github.com/jeanfric/embed"
)

var (
	imports = [...]string{}
)

const (
	decode string = `func(s string) (string, error) {
		return s, nil
	}`
)

func encode(contents io.Reader) (string, error) {
	b, err := ioutil.ReadAll(contents)
	if err != nil {
		return "", err
	}
	return strconv.Quote(string(b)), nil
}

// NewSequential creates a new sequential quoteembedder asset embedder.
func NewSequential() embed.AssetEmbedder {
	return embed.NewSequentialEmbedder(encode, decode, imports[:])
}

// NewConcurrent creates a new concurrent quoteembedder asset embedder.
func NewConcurrent() embed.AssetEmbedder {
	return embed.NewConcurrentEmbedder(encode, decode, imports[:])
}
