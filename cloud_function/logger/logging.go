package logger

import (
	"cloud.google.com/go/logging"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	"os"
)

// These variables are used for logging and are automatically
// set by the Cloud Functions runtime.
var (
	projectID    = os.Getenv("GCLOUD_PROJECT")
	functionName = os.Getenv("FUNCTION_NAME")
	region       = os.Getenv("FUNCTION_REGION")
)

// newLogger creates a Stackdriver logger.
func newLogger() (*logging.Client, *logging.Logger, error) {
	logClient, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("logging.NewClient: %v", err)
	}

	// Sets the ID of the log to write to.
	logID := "cloudfunctions.googleapis.com/cloud-functions"
	monitoredResource := monitoredres.MonitoredResource{
		Type: "cloud_function",
		Labels: map[string]string{
			"function_name": functionName,
			"region":        region,
		},
	}
	commonResource := logging.CommonResource(&monitoredResource)
	commonLabels := logging.CommonLabels(map[string]string{})
	logger := logClient.Logger(logID, commonResource, commonLabels)
	return logClient, logger, nil
}
