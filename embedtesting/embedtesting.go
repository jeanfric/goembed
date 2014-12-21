package embedtesting

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/jeanfric/goembed"
)

var (
	testAssets map[string]string = make(map[string]string)
)

func init() {
	var err error
	testAssets, err = loadAssets()
	if err != nil {
		panic(err)
	}
}

func GetTestAssets() []*goembed.Asset {
	return AssetsFromMap(testAssets)
}

func GetBenchAssets() []*goembed.Asset {
	benchAssets := make(map[string]string)

	// Amplify the size of the test data
	for k, v := range testAssets {
		for i := 0; i < 100; i++ {
			benchAssets[fmt.Sprintf("%i/%s", i, k)] = v
		}
	}

	return AssetsFromMap(benchAssets)
}

func AssetsFromMap(m map[string]string) []*goembed.Asset {
	assetList := make([]*goembed.Asset, 0, 0)
	for k, v := range m {
		assetList = append(assetList, &goembed.Asset{
			Reader: strings.NewReader(v),
			Key:    k,
		})
	}
	return assetList
}

func BenchmarkEmbedder(b *testing.B, ae goembed.AssetEmbedder) {
	assets := GetBenchAssets()
	var totBytes int64 = 0

	for i := 0; i < b.N; i++ {
		tf, err := ioutil.TempFile(os.TempDir(), "test")
		if err != nil {
			b.Fatal(err)
		}
		n, err := ae.AssetEmbed(tf, assets, "testing", "loadAssets")
		if err != nil {
			b.Fatal(err)
		}
		totBytes += int64(n)
		fName := tf.Name()
		tf.Close()
		err = os.Remove(fName)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(totBytes)
}
