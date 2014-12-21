package goembed

import (
	"io"
	"runtime"
)

// A concurrent embedder is an embedder that concurrently encodes
// assets, with up to runtime.NumCPU() concurrent embedders.
type ConcurrentEmbedder struct {
	encodeFunc func(contents io.Reader) (string, error)
	decodeFunc string
	imports    []string
}

// NewConcurrentEmbedder creates a new concurrent embedder that
// encodes assets using encodeFunc, with up to runtime.NumCPU()
// concurrent embedders.  The encoded assets are written to the
// generated Go source file, together with the decodeFunc and package
// imports.  The decodeFunc will wrap each string representation
// produced by the encodeFunc.
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
func NewConcurrentEmbedder(encodeFunc func(io.Reader) (string, error), decodeFunc string, imports []string) *ConcurrentEmbedder {
	return &ConcurrentEmbedder{
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
func (a *ConcurrentEmbedder) AssetEmbed(dst io.Writer, assets []*Asset, packageName, funcName string) (int, error) {
	assetChannel := make(chan *Asset, runtime.NumCPU())
	complete := make(chan *processedAsset, len(assets))
	for i := 0; i < runtime.NumCPU(); i++ {
		processAssets := func(assetChannel chan *Asset, complete chan *processedAsset) {
			for {
				select {
				case req, ok := <-assetChannel:
					if !ok {
						return
					}
					s, err := a.encodeFunc(req)
					complete <- &processedAsset{
						Asset: req,
						EncodedRepresentation: s,
						Error: err,
					}
				}
			}
		}
		go processAssets(assetChannel, complete)
	}

	for _, a := range assets {
		assetChannel <- a
	}

	queueResults := make(map[string]*processedAsset)
	for range assets {
		r := <-complete
		if r.Error != nil {
			return 0, r.Error
		}
		queueResults[r.Asset.Key] = r
	}
	close(assetChannel)
	close(complete)

	g := &generatedFileData{
		PackageName: packageName,
		FuncName:    funcName,
		Imports:     a.imports,
		DecodeFunc:  a.decodeFunc,
		Assets:      make([]*processedAsset, len(assets)),
	}
	for i, a := range assets {
		g.Assets[i] = queueResults[a.Key]
	}

	n, err := generateEmbedFile(dst, g)
	if err != nil {
		return n, err
	}
	return n, nil
}
