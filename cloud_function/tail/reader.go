package tail

import (
	"github.com/klauspost/compress/gzip"
	"io"
	"strings"
)


func useCompressReaderIfNeeded(reader io.ReadCloser, URL string) (io.ReadCloser, error) {
	if strings.HasSuffix(URL, ".gz") {
		return  gzip.NewReader(reader)
	}
	return reader, nil
}