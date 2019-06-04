package gcetoggle

import (
	"context"
	"errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

//InstanceRequest
type InstanceRequest struct {
	ProjectID    string
	Zone         string
	Name         string
	InstanceID   string
	TimeoutInSec int
	Action       string
	Selector     string
}

type StartRequest struct {
	InstanceRequest
}

type InstanceResponse struct {
	Error    string
	Status   string
	Instance *compute.Instance
}

type StartResponse InstanceResponse

type StopRequest struct {
	InstanceRequest
}


func (r *InstanceRequest) Init() {
	if r.ProjectID == "" {
		if credentials, err := google.FindDefaultCredentials(context.Background()); err == nil {
			r.ProjectID = credentials.ProjectID
		}
	}
	if r.TimeoutInSec == 0 {
		r.TimeoutInSec = DefaultTimeoutInSec
	}

}

func (r *InstanceRequest) Validate() error {
	if r.ProjectID == "" {
		return errors.New("projectID was empty")
	}
	if r.Zone == "" {
		return errors.New("zone was empty")
	}
	return nil

}

func NewInstanceResponse() *InstanceResponse {
	return &InstanceResponse{Status: StatusOk}
}

