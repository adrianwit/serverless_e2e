package cloud_function

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adrianwit/serverless_e2e/cloud_function/bq"
	"net/http"
	"strings"
	"time"
)

const maxExecutionTime = 60 * time.Second

//QueryRequest represents query request
type QueryRequest struct {
	ProjectID string
	DatasetID string
	UseLegacy bool
	SQL       string
}

//Validate checks if request is valid
func (r *QueryRequest) Validate() error {
	if r.SQL == "" {
		return fmt.Errorf("SQL was empty")
	}
	return nil
}

//QueryResponse represnets query response
type QueryResponse struct {
	Rows []map[string]bigquery.Value
}

//QueryFn runs SQL
func QueryFn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), maxExecutionTime)
	defer cancel()
	if strings.ToUpper(r.Method) != "POST" {
		NotImplemented(w, r)
		return
	}
	request := &QueryRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request: %v", err), http.StatusInternalServerError)
	}
	resonse, err := query(ctx, request)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to query: %v  %v", request.SQL, err), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(resonse)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to query: %v  %v", request.SQL, err), http.StatusInternalServerError)
	}
}

func query(ctx context.Context, request *QueryRequest) (*QueryResponse, error) {
	var response = &QueryResponse{
		Rows: make([]map[string]bigquery.Value, 0),
	}
	mapProvider := func() interface{} {
		var result = map[string]bigquery.Value{}
		return &result
	}
	err := bq.RunQuery(ctx, request.ProjectID, request.DatasetID, request.SQL, nil, request.UseLegacy, mapProvider, func(row interface{}) (b bool, e error) {
		aMap, ok := row.(*map[string]bigquery.Value)
		if !ok {
			return false, fmt.Errorf("expected *map[string]bigquery.Value, but had %T", aMap)
		}
		response.Rows = append(response.Rows, *aMap)
		return true, nil
	})
	return response, err
}

// NotImplemented returns status code 501 Not Implemented.
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Method %v is not yet implemented.", r.Method), http.StatusNotImplemented)
}
