package dispatcher

import (
	"context"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/cred"
	"os"
	"strings"
)

const slackCredentialsEnvKey = "slackSecrets"

var slackClient *slack.Client

func getClient() (*slack.Client, error) {
	if slackClient != nil {
		return slackClient, nil
	}
	credentials := os.Getenv(slackCredentialsEnvKey)
	if credentials == "" {
		return nil, fmt.Errorf("os.env.%v was empty", slackCredentialsEnvKey)
	}
	credConfig := cred.Config{}
	if err := credConfig.LoadFromReader(strings.NewReader(credentials), ".json"); err != nil {
		return nil, err
	}
	slackClient = slack.New(credConfig.Password)
	return slackClient, nil
}

func Post(ctx context.Context, target *Target, event interface{}, eventType string) error {
	client, err := getClient()
	JSON, err := toolbox.AsIndentJSONText(event)
	if err != nil {
		return err
	}
	request := slack.FileUploadParameters{
		Filename: eventType + ".json",
		Title:    eventType,
		Filetype: "json",
		Content:  string(JSON),
		Channels: []string{target.URL},
	}
	_, err = client.UploadFile(request)
	return err
}
