// Package goembed implements common asset embedder patterns that can
// be used by concrete asset embedder implementations.
package goembed

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/jeanfric/goembed/countingwriter"
)

// An Asset represents a named piece of data, typically the contents
// of a file identified by its path (beginning with a "/", and using
// forward slashes ("/") as path separators).
type Asset struct {
	io.Reader
	Key string
}

// A processedAsset represents an asset that has been encoded to a
// representation suitable for embedding in a Go source file.
type processedAsset struct {
	*Asset
	EncodedRepresentation string
	Error                 error
}

// AssetEmbedder is an interface that wraps the basic AssetEmbed method.
//
// AssetEmbed is a method that can produce a Go source file that
// embeds assets as encoded strings, in package packageName, and with
// fnName as the method to call to load and retrieve the embedded
// assets.  The fnName function must have the following signature:
//
//	func fnName() (map[string]string, error)
//
// The returned map is keyed by the asset key, and the value is the
// decoded contents of the asset (that is, an exact replica of the
// original contents of the piece of data that was embedded).
// Typically, the key is a file path, and the value is the contents of
// the file.
type AssetEmbedder interface {
	AssetEmbed(dst io.Writer, assets []*Asset, packageName, fnName string) (bytes int, err error)
}

// The generatedFileData structure contains all the information needed
// to produce a Go source file from a set of processed assets,
// complete with information about the decoding function, loading
// function, package name and list of imports.
type generatedFileData struct {
	PackageName string // The Go package name to use
	FuncName    string // The name of the loading function
	Imports     []string
	Assets      []*processedAsset
	DecodeFunc  string
}

// FindAssets walks a directory recursively and generates a list of
// embeddable assets that can be embedded using an AssetEmbedder.  The
// Key of each asset will start with a forward slash ("/"), and use
// slashes as path separators.
func FindAssets(rootPath string) ([]*Asset, error) {
	assetList := make([]*Asset, 0, 0)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// We could just pass along the opened file,
		// but let's just be done with the reading
		// here.  This way, we can exit early if there
		// are issues reading some of the files.
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		assetList = append(assetList, &Asset{
			Reader: bytes.NewReader(b),
			Key:    "/" + filepath.ToSlash(relPath),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return assetList, nil
}

func generateEmbedFile(dst io.Writer, data *generatedFileData) (int, error) {
	// TODO: using templates is probably a tad overkill here, but
	// it makes the code more pleasant to read.
	outputTemplate := `package {{.PackageName}}
`

	if len(data.Imports) > 0 {
		outputTemplate += `
import ({{range $v := .Imports}}
	{{printf "%q" $v}}{{end}}
)
`
	}

	outputTemplate += `
func {{.FuncName}}() (map[string]string, error) {
	decode := {{.DecodeFunc}}

	var a string
	var err error
	assets := make(map[string]string)
{{range $i, $v := .Assets}}
	a, err = decode({{$v.EncodedRepresentation}})
	if err != nil {
		return nil, err
	}
	assets[{{printf "%q" $v.Key}}] = a
{{end}}
	return assets, nil
}
`
	t := template.Must(template.New("").Parse(outputTemplate))

	// The counting writer will enable us to report how many bytes
	// we have written, since t.Execute does not provide this
	// information.
	countingWriter := countingwriter.New(dst)
	err := t.Execute(countingWriter, data)
	if err != nil {
		return countingWriter.BytesWritten(), err
	}
	return countingWriter.BytesWritten(), nil
}
