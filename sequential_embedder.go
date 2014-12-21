package goembed

import (
	"io"
)

// A sequential embedder is an embedder that encodes assets one by one.
type SequentialEmbedder struct {
	encodeFunc func(contents io.Reader) (string, error)
	decodeFunc string
	imports    []string
}

// NewSequentialEmbedder creates a new sequential embedder that
// encodes assets using encodeFunc.  The encoded assets are written to
// the generated Go source file, together with the decodeFunc and
// package imports.  The decodeFunc will wrap each string
// representation produced by the encodeFunc.
//
// The encodeFunc should return a string enclosed in its delimiters
// (double quotes or backticks).  The embedder implementation can thus
// choose if it wants to return a quoted string (with double quotes)
// or a raw string (with backticks).
//
// The decodeFunc must be written in this form:
//
//	func(s string) (string, error) {
//		// ...
//	}
//
// The argument is the encoded string, and the returned string is the
// decoded data (matching the original asset data).
//
// The imports array should match the compilation requirements of the
// decodeFunc.
func NewSequentialEmbedder(encodeFunc func(io.Reader) (string, error), decodeFunc string, imports []string) *SequentialEmbedder {
	return &SequentialEmbedder{
		encodeFunc: encodeFunc,
		decodeFunc: decodeFunc,
		imports:    imports,
	}
}

// AssetEmbed outputs a Go source file containing the assets.  The
// source file will be in package packageName, and the function that
// returns the assets will be named funcName.  This function will have
// the following signature:
//
// 	func funcName() (map[string]string, error)
func (e *SequentialEmbedder) AssetEmbed(dst io.Writer, assets []*Asset, packageName, funcName string) (int, error) {
	g := &generatedFileData{
		PackageName: packageName,
		FuncName:    funcName,
		Imports:     e.imports,
		DecodeFunc:  e.decodeFunc,
		Assets:      make([]*processedAsset, len(assets)),
	}
	for i, a := range assets {
		r, err := e.encodeFunc(a)
		if err != nil {
			return 0, err
		}

		g.Assets[i] = &processedAsset{
			Asset: a,
			EncodedRepresentation: r,
		}
	}

	n, err := generateEmbedFile(dst, g)
	if err != nil {
		return n, err
	}
	return n, nil
}
