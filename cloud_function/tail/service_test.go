package tail

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/viant/dsunit"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/url"
	"golang.org/x/net/context"
	"log"
	"path"
	"testing"
)

func TestService_Transfer(t *testing.T) {
	parent := toolbox.CallerDirectory(3)
	if !dsunit.InitFromURL(t, path.Join(parent, "test", "config.yaml")) {
		return
	}
	config, err := NewConfigFromURL(url.NewResource(path.Join(parent, "test/config.json")).URL)
	if !assert.Nil(t, err) {
		log.Fatal(err)
	}
	service, err := New(config)
	if !assert.Nil(t, err) {
		log.Fatal(err)
	}

	useCases := []struct {
		description string
		resourceURL string
		expectPath  string
		transferred int
		hasError    bool
		hasMeta     bool
	}{
		{
			description: "basic transfer",
			resourceURL: url.NewResource(path.Join(parent, "test/data1/events/data_8D05F5F23F76.json")).URL,
			expectPath:  path.Join(parent, "test/data1/expect"),
			transferred: 11,
		},
	}

	for _, useCase := range useCases {
		_, name := toolbox.URLSplit(useCase.resourceURL)
		err = service.Transfer(context.Background(), useCase.resourceURL)
		meta, err := NewMetaFromURL(toolbox.URLPathJoin(config.ProcessedURL, name+".meta"))
		if useCase.hasMeta {
			assert.NotNil(t, err, useCase.description)
		}
		if useCase.hasError {
			assert.NotNil(t, err, useCase.description)
			if useCase.hasMeta {
				assert.True(t, meta.Error != "", useCase.description)
			}
			continue
		}
		if !assert.Nil(t, err) {
			log.Print(err)
			continue
		}
		assert.Equal(t, useCase.transferred, meta.Transferred)
		expectedData := dsunit.NewDatasetResource("db1", useCase.expectPath, "", "")
		dsunit.Expect(t, dsunit.NewExpectRequest(dsunit.FullTableDatasetCheckPolicy, expectedData))
	}

}
