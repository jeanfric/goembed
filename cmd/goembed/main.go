// Goembed generates a file named "assets.generated.go" containing an
// encoded version of the contents of a directory.
//
// The generated source code has no external import dependencies; it
// relies solely on the standard Go library.
//
// The generated code provides a function that loads and returns the
// embedded assets, with the following signature:
//
// 	func loadAssets() (map[string]string, error)
//
// Given this "testdata" directory:
//
//	testdata/
//	|- index.html
//	`- img
//	   `gopher.png
//
// Then, after running "goembed testdata" and compiling the package,
// calling loadAssets will return a map like the following:
//
//	m["/index.html"] = "<html><head>..."
//	m["/img/gopher.png"] = "\x89PNG\r\n\x1a\n..."
//
// The paths will all begin at "/" and use forward slashes ("/") as
// path separators.
//
// Goembed is useful in combination with "go generate" to bundle static
// assets in a program binary.  For example, to embed all files under the
// "static" directory:
//
// 	package main
// 	//go:generate embed static
// 	[...]
//
// 	$ go generate
// 	$ go build
//
// Goembed supports encoding the data using the following algorithms:
//
//	* quote: quoted Go string
//	* hex: hex-encoded
//	* base64: base64-encoded
//	* zhex: zlib-compressed, hex-encoded
//	* zbase64: zlib-compressed, base64-encoded
//
// Usage:
//	goembed [-package p] [-func f] [-o output] directory
//
// The flags and their default values are:
//	-c=false
//		use concurrent version of the chosen algorithm
//	-e="quote"
//		embedding algorithm
//	-func="loadAssets"
//		name of loading function
//	-o="assets.generated.go"
//		name of generated file
//	-package="main"
//		package of the generated source file (if $GOPACKAGE is
//		set, such as when using "go generate", $GOPACKAGE
//		takes precedence)
//
// See also: package github.com/jeanfric/embedfs implements an
// http.FileSystem backed by a map[string]string, compatible directly
// with the map returned by loadAssets.  See
// https://github.com/jeanfric/embedfs for more information.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/jeanfric/goembed"
	"github.com/jeanfric/goembed/base64embedder"
	"github.com/jeanfric/goembed/hexembedder"
	"github.com/jeanfric/goembed/quoteembedder"
	"github.com/jeanfric/goembed/zbase64embedder"
	"github.com/jeanfric/goembed/zhexembedder"
)

func usage() {
	details := `
usage: goembed [-package p] [-func f] [-o output] directory

Goembed generates a file named "assets.generated.go" containing an
encoded version of the contents of the specified directory.

The generated file contains a function that loads and returns the
embedded assets, with the following signature:

	func loadAssets() (map[string]string, error)

Goembed is useful in combination with "go generate" to bundle static
assets in a program binary.  For example, to embed all files under the
"static" directory:

	package main
	//go:generate goembed static
	[...]

	$ go generate
	$ go build

See also: https://github.com/jeanfric/embedfs.

Goembed flags:
`

	fmt.Fprintf(os.Stderr, details)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var destFile, packageName, fnName, embedder string
	var concurrent bool
	flag.StringVar(&packageName, "package", "main", "package of the generated source file (if $GOPACKAGE is set, such as when using \"go generate\", $GOPACKAGE takes precedence)")
	flag.StringVar(&fnName, "func", "loadAssets", "name of loading function")
	flag.StringVar(&destFile, "o", "assets.generated.go", "name of generated file")
	flag.StringVar(&embedder, "e", "quote", "embedding algorithm")
	flag.BoolVar(&concurrent, "c", false, "use concurrent version of the chosen algorithm")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
	}

	srcPath := flag.Args()[0]

	// Go generate will pass us the package name of the file
	// containing the go:generate directive.  Use that by default.
	envPackage := os.Getenv("GOPACKAGE")
	if envPackage != "" {
		packageName = envPackage
	}

	dest, err := os.Create(destFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	assets, err := goembed.FindAssets(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	var ae goembed.AssetEmbedder

	switch embedder {
	case "zbase64":
		if concurrent {
			ae = zbase64embedder.NewConcurrent()
		} else {
			ae = zbase64embedder.NewSequential()
		}
	case "base64":
		if concurrent {
			ae = base64embedder.NewConcurrent()
		} else {
			ae = base64embedder.NewSequential()
		}
	case "hex":
		if concurrent {
			ae = hexembedder.NewConcurrent()
		} else {
			ae = hexembedder.NewSequential()
		}
	case "zhex":
		if concurrent {
			ae = zhexembedder.NewConcurrent()
		} else {
			ae = zhexembedder.NewSequential()
		}
	case "quote":
		if concurrent {
			ae = quoteembedder.NewConcurrent()
		} else {
			ae = quoteembedder.NewSequential()
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown embedding algorithm \"%s\"\n", embedder)
		os.Exit(1)
	}

	if _, err := ae.AssetEmbed(dest, assets, packageName, fnName); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	return
}
