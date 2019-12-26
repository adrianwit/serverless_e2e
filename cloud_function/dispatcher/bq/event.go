package bq

import "google.golang.org/api/bigquery/v2"

//Event represents dispatcher job event
type Event bigquery.Job

//Type returns job type
func (e *Event) Type() string {
	return e.Configuration.JobType
}

//Source returns job source
func (e *Event) Source() string {
	switch e.Configuration.JobType {
	case "QUERY":
		return e.Configuration.Query.Query
	case "LOAD":
		return ""
	case "EXTRACT":
		source := e.Configuration.Extract.SourceTable
		return source.ProjectId + ":" + source.DatasetId + "." + source.TableId
	case "COPY":
		source := e.Configuration.Copy.SourceTable
		return source.ProjectId + ":" + source.DatasetId + "." + source.TableId
	}
	return ""
}

//Destination returns job destination
func (e *Event) Destination() string {
	switch e.Configuration.JobType {
	case "QUERY":
		dest := e.Configuration.Query.DestinationTable
		if dest == nil {
			return ""
		}
		return dest.ProjectId + ":" + dest.DatasetId + "." + dest.TableId
	case "LOAD":
		dest := e.Configuration.Load.DestinationTable
		if dest == nil {
			return ""
		}
		return dest.ProjectId + ":" + dest.DatasetId + "." + dest.TableId
	case "EXTRACT":
		return e.Configuration.Extract.DestinationUri
	case "COPY":
		source := e.Configuration.Copy.DestinationTable
		return source.ProjectId + ":" + source.DatasetId + "." + source.TableId
	}
	return ""
}
