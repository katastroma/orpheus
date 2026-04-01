//revive:disable:package-comments
package helm

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

const directoryPrefix = "chart/"

// repackage reads a plain tar stream and produces a gzipped tar stream with
// every entry nested under a directory prefix. See Render for why this exists.
func repackage(r io.Reader) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		gw := gzip.NewWriter(pw)
		tw := tar.NewWriter(gw)
		tr := tar.NewReader(r)

		err := rewriteEntries(tr, tw)

		if closeErr := tw.Close(); err == nil {
			err = closeErr
		}
		if closeErr := gw.Close(); err == nil {
			err = closeErr
		}

		pw.CloseWithError(err)
	}()

	return pr
}

func rewriteEntries(tr *tar.Reader, tw *tar.Writer) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		hdr.Name = directoryPrefix + hdr.Name

		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err = io.Copy(tw, tr); err != nil {
			return err
		}
	}
}
