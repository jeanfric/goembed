package base64embedder

import (
	"testing"

	"github.com/jeanfric/goembed/embedtesting"
)

func BenchmarkSequentialEmbedder(b *testing.B) {
	embedtesting.BenchmarkEmbedder(b, NewSequential())
}

func BenchmarkConcurrentEmbedder(b *testing.B) {
	embedtesting.BenchmarkEmbedder(b, NewConcurrent())
}
