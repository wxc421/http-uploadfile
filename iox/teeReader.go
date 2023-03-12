package iox

import (
	"io"
)

// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the writer must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func TeeReader(r io.ReadCloser, w io.Writer) io.ReadCloser {
	return &teeReader{r, w}
}

type teeReader struct {
	r io.ReadCloser
	w io.Writer
}

func (t *teeReader) Close() error {
	// todo export status
	return t.r.Close()
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err = t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}
