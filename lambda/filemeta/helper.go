package filemeta

import (
	"compress/gzip"
	"github.com/viant/toolbox/cred"
	"github.com/viant/toolbox/storage"
	"github.com/viant/toolbox/storage/s3"
	"io/ioutil"
	"path"
	"strings"
	"unsafe"
)

func countObjectLines(service storage.Service, URL string) (int, error) {
	isCompressed := path.Ext(URL) == ".gz"
	contentReader, err := service.DownloadWithURL(URL)
	if err != nil {
		return 0, err
	}
	defer contentReader.Close()
	reader := contentReader
	if isCompressed {
		if reader, err = gzip.NewReader(contentReader); err != nil {
			return 0, err
		}
	}
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return 0, err
	}
	fragment := *(*string)(unsafe.Pointer(&content))
	lineCount := strings.Count(fragment, "\n")
	return lineCount, nil
}

func getStorageService(region string) (storage.Service, error) {
	s3.SetProvider(&cred.Config{Region: region})
	return storage.NewServiceForURL("s3://bucket/", "")
}
