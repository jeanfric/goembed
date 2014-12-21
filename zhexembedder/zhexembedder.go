// Package zhexembedder implements an asset embedder that compresses
// assets using zlib, then encodes the resulting data as hexadecimal
// strings.
package zhexembedder

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"io"

	"github.com/jeanfric/goembed"
)

var (
	imports = [...]string{"bytes", "compress/zlib", "encoding/hex", "io/ioutil"}
)

const (
	decode = `func(s string) (string, error) {
		b, err := hex.DecodeString(s)
		if err != nil {
			return "", err
		}
		r, err := zlib.NewReader(bytes.NewReader(b))
		if err != nil {
			return "", err
		}
		defer r.Close()
		ob, err := ioutil.ReadAll(r)
		if err != nil {
			return "", err
		}
		return string(ob), nil
	}`
)

func encode(contents io.Reader) (string, error) {
	var zb bytes.Buffer
	w := zlib.NewWriter(&zb)
	_, err := io.Copy(w, contents)
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	return "`" + hex.EncodeToString(zb.Bytes()) + "`", nil
}

// NewSequential creates a new sequential zhexembedder asset embedder.
func NewSequential() goembed.AssetEmbedder {
	return goembed.NewSequentialEmbedder(encode, decode, imports[:])
}

// NewConcurrent creates a new concurrent zhexembedder asset embedder.
func NewConcurrent() goembed.AssetEmbedder {
	return goembed.NewConcurrentEmbedder(encode, decode, imports[:])
}
