// Package countingwriter implements an io.Writer that wraps another
// writer and keeps track of the number of bytes that are written by
// that wrapped writer.
package countingwriter

import "io"

type CountingWriter struct {
	writer       io.Writer
	bytesWritten int
}

// New creates a new writer that wraps w.  The wrapping writer counts
// the number of bytes written to the wrapped writer.
func New(w io.Writer) *CountingWriter {
	return &CountingWriter{
		writer:       w,
		bytesWritten: 0,
	}
}

func (w *CountingWriter) Write(b []byte) (int, error) {
	n, err := w.writer.Write(b)
	w.bytesWritten += n
	return n, err
}

// BytesWritten returns the number of bytes that were written to the wrapped writer.
func (w *CountingWriter) BytesWritten() int {
	return w.bytesWritten
}
