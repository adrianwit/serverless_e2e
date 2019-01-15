package loginfo

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"github.com/viant/toolbox/storage"
	"strings"
	"testing"
)

func Test_countLines(t *testing.T) {
	service := storage.NewMemoryService()
	var useCases = []struct {
		description   string
		compressed    bool
		content       []byte
		expectedCount int
	}{
		{
			description:   "text file",
			content:       []byte(strings.Repeat("abc\nz", 128)),
			expectedCount: 128,
		},

		{
			description:   "text file",
			compressed:    true,
			content:       []byte(strings.Repeat("xyz\ny", 1024)),
			expectedCount: 1024,
		},
	}

	for _, useCase := range useCases {
		content := useCase.content
		URL := "mem:///test.txt"
		if useCase.compressed { //compress if needed
			buf := new(bytes.Buffer)
			writer := gzip.NewWriter(buf)
			_, err := writer.Write(useCase.content)
			assert.Nil(t, err, useCase.description)
			_ = writer.Close()
			content = buf.Bytes()
			URL += ".gz"
		}
		err := service.Upload(URL, bytes.NewReader(content))
		assert.Nil(t, err, useCase.description)
		response := &Response{}
		object, err := service.StorageObject(URL)
		assert.Nil(t, err, useCase.description)
		countLines(service, object, response)
		assert.EqualValues(t, useCase.expectedCount, response.LinesCount)
	}

}

func Test_countFileAndLines(t *testing.T) {
	var useCases = []struct {
		description       string
		url               string
		assets            map[string][]byte
		expectedRowCount  int
		expectedFileCount int
	}{
		{
			description: "single folder count",
			assets: map[string][]byte{
				"mem:///folder1/asset1.txt": []byte(strings.Repeat("abc\nz", 128)),
				"mem:///folder1/asset2.txt": []byte(strings.Repeat("abc\nz", 64)),
			},
			url:               "mem:///folder1",
			expectedRowCount:  192,
			expectedFileCount: 2,
		},

		{
			description: "nester folder count",
			assets: map[string][]byte{
				"mem:///folder2/asset1.txt":     []byte(strings.Repeat("abc\nz", 128)),
				"mem:///folder2/asset2.txt":     []byte(strings.Repeat("abc\nz", 64)),
				"mem:///folder2/sub/asset3.txt": []byte(strings.Repeat("abc\nz", 64)),
			},
			url:               "mem:///folder2",
			expectedRowCount:  256,
			expectedFileCount: 3,
		},
	}

	for _, useCase := range useCases {
		service := storage.NewMemoryService()

		for url, asset := range useCase.assets {
			err := service.Upload(url, bytes.NewReader(asset))
			assert.Nil(t, err, useCase.description)
		}
		response := &Response{}
		countFileAndLines(service, useCase.url, response)
		assert.EqualValues(t, useCase.expectedRowCount, response.LinesCount, "line count,  " + useCase.description)
		assert.EqualValues(t, useCase.expectedFileCount, response.FileCount, "file count,  " + useCase.description)
		assert.EqualValues(t, "", response.Status, useCase.description)
	}
}
