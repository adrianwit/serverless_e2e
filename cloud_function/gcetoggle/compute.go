package gcetoggle

import (
	"context"
	"google.golang.org/api/compute/v1"
)


//NewComputeService returns new compute computeService
func NewComputeService(ctx context.Context) (*compute.Service, error){
	computeService, err := compute.NewService(ctx)
	return computeService, err
}