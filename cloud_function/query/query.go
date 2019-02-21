package query

import (
	"cloud.google.com/go/bigquery"
	"context"
	"google.golang.org/api/iterator"
	"os"
)

var client *bigquery.Client

func getClient(project string) (*bigquery.Client, error) {
	if client != nil {
		return client, nil
	}
	if project == "" {
		project = os.Getenv("GCLOUD_PROJECT")
	}
	var err error
	client, err = bigquery.NewClient(context.Background(), project)
	return client, err
}

func RunQuery(ctx context.Context, project, dataset string, SQL string, params []interface{}, useLegacy bool, rowProvider func() interface{}, handler func(row interface{}) (bool, error)) error {
	client, err := getClient(project)
	if err != nil {
		return err
	}
	query := client.Query(SQL)
	query.UseLegacySQL = useLegacy
	if len(params) > 0 {
		query.Parameters = make([]bigquery.QueryParameter, 0)
		for _, param := range params {
			query.Parameters = append(query.Parameters, bigquery.QueryParameter{Value: param})
		}
	}
	query.DefaultDatasetID = dataset
	query.DefaultProjectID = project
	job, err := query.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	it, err := job.Read(ctx)
	if err != nil {
		return err
	}
	for {
		row := rowProvider()
		err := it.Next(row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		cont, e := handler(row)
		if e != nil || !cont {
			return e
		}
	}
	return nil
}
