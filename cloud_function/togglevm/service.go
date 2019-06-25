package togglevm

import (
	"context"
	"fmt"
	"google.golang.org/api/compute/v1"
	"strconv"
	"time"
)

type Service interface {
	Start(*StartRequest) *InstanceResponse
	Stop(*StopRequest) *InstanceResponse
}

type service struct {
	*compute.Service
	Context context.Context
}

func (s *service) Start(request *StartRequest) *InstanceResponse {
	response := NewInstanceResponse()
	err := s.start(request, response)
	if err != nil {
		response.Error = err.Error()
		response.Status = StatusError
	}
	return response
}


func (s *service) Stop(request *StopRequest) *InstanceResponse {
	response := NewInstanceResponse()
	err := s.stop(request, response)
	if err != nil {
		response.Error = err.Error()
		response.Status = StatusError
	}
	return response
}


func (s *service) getInstance(instanceService *compute.InstancesService, request *InstanceRequest) (*compute.Instance, error) {
	var err error
	if request.InstanceID != "" {
		instanceCall := instanceService.Get(request.ProjectID, request.Zone, request.InstanceID)
		instanceCall.Context(s.Context)
		return instanceCall.Do()
	}
	pageToken := ""
	var instanceList *compute.InstanceList
	listCall := instanceService.List(request.ProjectID, request.Zone)
	for ; ; {

		listCall.PageToken(pageToken)
		listCall.Context(s.Context)
		instanceList, err = listCall.Do()
		if err != nil || len(instanceList.Items) == 0 {
			return nil, err
		}
		for i := range instanceList.Items {
			if instanceList.Items[i].Name == request.Name {
				request.InstanceID = strconv.FormatUint(instanceList.Items[i].Id, 10)
				return instanceList.Items[i], nil
			}
		}
		pageToken = instanceList.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return nil, nil
}



func (s *service) waitForStatus(instanceService *compute.InstancesService, request *InstanceRequest, status string) (*compute.Instance, error) {
	startTime := time.Now()
	maxWaitTime := time.Duration(request.TimeoutInSec) * time.Second
	for  {
		instance, err := s.getInstance(instanceService, request)
		if err != nil {
			return nil,  err
		}
		if instance.Status == status {
			return instance, nil
		}
		time.Sleep(3 * time.Second)
		if time.Now().Sub(startTime) > maxWaitTime {
			break
		}
	}
	return nil, fmt.Errorf("exceeded max wait time %s", maxWaitTime)
}


func (s *service) start(request *StartRequest, response *InstanceResponse) error {
	request.Init()
	if err := request.Validate();err != nil {
		return err
	}
	instanceService := compute.NewInstancesService(s.Service)
	instance, err := s.getInstance(instanceService, &request.InstanceRequest)
	if err != nil {
		return err
	}

	response.Instance = instance
	switch instance.Status {
	case StatusRunning:
		return nil
	case StatusTerminated:
		startCall := instanceService.Start(request.ProjectID, request.Zone, request.InstanceID)
		startCall.Context(s.Context)
		_, err := startCall.Do()
		if err != nil {
			return fmt.Errorf("failed to start instance %v: %v", request.InstanceID, err)
		}
		response.Instance, err = s.waitForStatus(instanceService, &request.InstanceRequest, StatusRunning)
	}
	return nil
}



func (s *service) stop(request *StopRequest, response *InstanceResponse) error {
	request.Init()
	if err := request.Validate();err != nil {
		return err
	}
	instanceService := compute.NewInstancesService(s.Service)
	instance, err := s.getInstance(instanceService, &request.InstanceRequest)
	if err != nil {
		return err
	}

	response.Instance = instance
	switch instance.Status {
	case StatusTerminated:
		return nil
	case StatusRunning:
		stopCall := instanceService.Stop(request.ProjectID, request.Zone, request.InstanceID)
		stopCall.Context(s.Context)
		_, err := stopCall.Do()
		if err != nil {
			return fmt.Errorf("failed to stop instance %v: %v", request.InstanceID, err)
		}
		fallthrough
	default:
		response.Instance, err = s.waitForStatus(instanceService, &request.InstanceRequest, StatusRunning)
	}
	return nil

}

func New() (Service, error) {
	var err error
	result := &service{Context: context.Background()}
	result.Service, err = NewComputeService(result.Context)
	return result, err
}
