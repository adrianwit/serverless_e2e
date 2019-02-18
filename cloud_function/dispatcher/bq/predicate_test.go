package bq

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/bigquery/v2"
	"testing"
)

func TestPredicate_Apply(t *testing.T) {

	var useCases = []struct {
		description string
		filter      interface{}
		event       interface{}
		hasError    bool
		expect      bool
	}{

		{
			description: "filter by source dataset match",
			expect:      true,
			filter: &Filter{
				Source: ".+:myDataset\\..+",
			},
			event: &Event{
				Configuration: &bigquery.JobConfiguration{
					JobType: "COPY",
					Copy: &bigquery.JobConfigurationTableCopy{
						SourceTable: &bigquery.TableReference{
							ProjectId: "myProject",
							DatasetId: "myDataset",
							TableId:   "sourceTable",
						},
						DestinationTable: &bigquery.TableReference{
							ProjectId: "myProject",
							DatasetId: "myDataset",
							TableId:   "myDestTable",
						},
					},
				},
			},
		},
		{
			description: "filter by target table does match",
			expect:      false,
			filter: &Filter{
				Destination: ".+myDestTable2",
			},
			event: &Event{
				Configuration: &bigquery.JobConfiguration{
					JobType: "LOAD",
					Load: &bigquery.JobConfigurationLoad{
						DestinationTable: &bigquery.TableReference{
							ProjectId: "myProject",
							DatasetId: "myDataset",
							TableId:   "myDestTable1",
						},
					},
				},
			},
		},

		{
			description: "filter by job type and URI match",
			expect:      true,
			filter: &Filter{
				Type:        "EXTRACT",
				Destination: ".+table10",
			},
			event: &Event{
				Configuration: &bigquery.JobConfiguration{
					JobType: "EXTRACT",
					Extract: &bigquery.JobConfigurationExtract{
						SourceTable: &bigquery.TableReference{
							ProjectId: "myProject",
							DatasetId: "myDataset",
							TableId:   "mySourceTable",
						},
						DestinationUri: "gs://myBucket/table10",
					},
				},
			},
		},
		{
			description: "filter by job type and URI does not match",
			expect:      false,
			filter: &Filter{
				Type:        "EXTRACT",
				Destination: ".+table2",
			},
			event: &Event{
				Configuration: &bigquery.JobConfiguration{
					JobType: "EXTRACT",
					Extract: &bigquery.JobConfigurationExtract{
						SourceTable: &bigquery.TableReference{
							ProjectId: "myProject",
							DatasetId: "myDataset",
							TableId:   "mySourceTable",
						},
						DestinationUri: "gs://myBucket/table10",
					},
				},
			},
		},

		{
			description: "filter by source match",
			expect:      true,
			filter: &Filter{
				Type:   "QUERY",
				Source: ".+tableX",
			},
			event: &Event{
				Configuration: &bigquery.JobConfiguration{
					JobType: "QUERY",
					Query: &bigquery.JobConfigurationQuery{
						DestinationTable: &bigquery.TableReference{
							ProjectId: "myProject",
							DatasetId: "myDataset",
							TableId:   "mySourceTable",
						},
						Query: "SELECT  FROM tableX WHERE 1 = 0",
					},
				},
			},
		},

		{
			description: "filter by source does not match",
			expect:      false,
			filter: &Filter{
				Type:   "QUERY",
				Source: ".+tableX.*",
			},
			event: &Event{
				Configuration: &bigquery.JobConfiguration{
					JobType: "QUERY",
					Query: &bigquery.JobConfigurationQuery{
						DestinationTable: &bigquery.TableReference{
							ProjectId: "myProject",
							DatasetId: "myDataset",
							TableId:   "mySourceTable",
						},
						Query: "SELECT  FROM tableY",
					},
				},
			},
		},
		{
			description: "invalid filter error",
			hasError:    true,
		},
		{
			description: "invalid event type",
			filter:      map[string]interface{}{},
			event:       map[string]interface{}{},
			expect:      false,
		},
	}

	for _, useCase := range useCases {
		predicate, err := NewPredicate(useCase.filter)
		if useCase.hasError {
			assert.NotNil(t, err, useCase.description)
			continue
		}
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		actual := predicate.Apply(useCase.event)
		assert.Equal(t, useCase.expect, actual, useCase.description)
	}
}
