package split

import (
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/storage"
	"github.com/viant/toolbox/storage/s3"
	_ "github.com/viant/toolbox/storage/gs"
	"net/url"
	"golang.org/x/net/context"
)


const gsResourceURL = "gs://dummy/dummy"

type Service interface {
	Mirror(SourceURL string) error
}


type service struct {
	config *Config
	destService storage.Service
	sourceService storage.Service
}

//Mirror mirrors source with destination
func (s *service) Mirror(SourceURL string) error {
	reader, err := s.sourceService.DownloadWithURL(SourceURL)
	if err != nil {
		return err
	}
	defer reader.Close()
	_, name := toolbox.URLSplit(SourceURL)
	destURL := toolbox.URLPathJoin(s.config.DestURL, name)
	return  s.destService.Upload(destURL, reader)
}

func newService(ctx context.Context, config *Config) (Service, error) {
	credConfig, err := getSecret(ctx, config.SecretKey, config.Secrets)
	if err != nil {
		return nil, err
	}
	URL, err := url.Parse(config.DestURL)
	if err != nil {
		return nil, err
	}
	switch URL.Scheme {
		case "s3":
			s3.SetProvider(credConfig)
	}
	result := &service{
		config:config,
	}
	result.destService, err = storage.NewServiceForURL(config.DestURL, "")
	if err != nil {
		return nil, err
	}
	result.sourceService, err = storage.NewServiceForURL(gsResourceURL, "")
	if err != nil {
		return nil, err
	}
	return result, nil
}






var srv Service

//GetService returns service
func GetService(ctx context.Context, config *Config) (Service, error) {
	if srv != nil {
		return srv, nil
	}
	return newService(ctx, config)
}
