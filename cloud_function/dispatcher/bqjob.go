package dispatcher

import (
	"golang.org/x/net/context"
	"google.golang.org/api/bigquery/v2"
)

func GetBQJob(ctx context.Context, projectID, jobID string) (*bigquery.Job, error) {
	httpClient, err := getDefaultClient(ctx, bigquery.CloudPlatformScope, bigquery.BigqueryScope, bigquery.BigqueryInsertdataScope)
	if err != nil {
		return nil, err
	}
	service, err := bigquery.New(httpClient)
	if err != nil {
		return nil, err
	}
	jobService := bigquery.NewJobsService(service)
	call := jobService.Get(projectID, jobID)
	call.Context(ctx)
	return call.Do()
}
