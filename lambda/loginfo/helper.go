package loginfo

import (
	"compress/gzip"
	"github.com/viant/toolbox/storage"
	"io/ioutil"
	"path"
	"strings"
	"sync/atomic"
	"unsafe"
)



func countLines(service storage.Service, object storage.Object, response *Response) {
	atomic.AddUint32(&response.FileCount, 1)
	isCompressed := path.Ext(object.URL()) == ".gz"
	contentReader, err := service.Download(object)
	if response.SetError(err) {
		return
	}
	defer contentReader.Close()
	reader := contentReader
	if isCompressed {
		if reader, err = gzip.NewReader(contentReader); err != nil {
			response.SetError(err)
			return
		}
	}
	content, err := ioutil.ReadAll(reader)
	if response.SetError(err) {
		return
	}
	fragment := *(*string)(unsafe.Pointer(&content))
	lineCount := strings.Count(fragment, "\n")
	atomic.AddUint32(&response.LinesCount, uint32(lineCount))
}



func countFileAndLines(service storage.Service, URL string, response *Response) {
	objects, err := service.List(URL)
	if response.SetError(err) {
		return
	}
	URL = strings.Trim(URL, "/")
	for _, object := range objects {
		if strings.Trim(object.URL(), "/") ==  URL {
			continue
		}
		if object.IsFolder() {
			countFileAndLines(service, object.URL(), response)
			continue
		}
		countLines(service, object, response)
	}
}
