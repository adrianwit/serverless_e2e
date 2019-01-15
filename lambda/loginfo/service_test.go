package loginfo

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/toolbox"
	"os"
	"path"
	"testing"
)

func TestService_CountLogs(t *testing.T) {
	credentialsFile := path.Join(os.Getenv("HOME"), ".secret/aws.json")
	if ! toolbox.FileExists(credentialsFile) {
		return
	}
	service := New()
	response := service.CountLogs(&Request{
		Region:"us-east-2",
		URL:"s3://mye2e-bucket/folder1/",
	})
	assert.EqualValues(t, 2, response.FileCount)
}
