// The goembedtest command tests that goembed-embedded assets match
// with the original source files that were embedded.
//
// To use this command, emebedded assets must first be generated from
// a given directory:
//
//	$ goembed testdata
//
// Then, embedded assets must be compiled into goembedtest:
//
//	$ go build
//
// When executed, goembedtest will then compare its embedded assets
// with the original files in the directory specified:
//
//	$ goembedtest testdata
//
// If the assets don't match with the original files, goembedtest will
// exit with a failure status code.
package main

//go:generate goembed testdata

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: goembedtest [-q] directory\n")
	flag.PrintDefaults()
	os.Exit(2)
}
func main() {
	var quiet bool
	flag.BoolVar(&quiet, "q", false, "quiet")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
	}

	rootPath := flag.Args()[0]

	assets, err := loadAssets()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	var results []string
	fail := false

	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
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
		fileData, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		assetData := []byte(assets["/"+filepath.ToSlash(relPath)])

		if bytes.Compare(fileData, assetData) != 0 {
			results = append(results, fmt.Sprintf("fail\t\"%v\" (%v bytes) != asset[\"%v\"] (%v bytes)", path, len(fileData), relPath, len(assetData)))
			fail = true
		} else {
			results = append(results, fmt.Sprintf("ok\t\"%v\" == asset[\"%v\"] (%v bytes)", path, relPath, len(assetData)))
		}

		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if !quiet {
		if fail {
			fmt.Fprintf(os.Stdout, "FAIL\n")
			for _, v := range results {
				fmt.Fprintf(os.Stdout, "%v\n", v)
			}
		} else {
			fmt.Fprintf(os.Stdout, "PASS\t%v comparisons\n", len(results))

		}
	}
	if fail {
		os.Exit(1)
	}

}
