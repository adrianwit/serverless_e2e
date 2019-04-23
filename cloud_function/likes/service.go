package cloud_function

import (
	"context"
	"firebase.google.com/go/db"
)

type Service interface {
	UpdateLikes(ctx context.Context, resource Resource) error
}

type service struct{}

func (s *service) update(snapshot db.TransactionNode) (interface{}, error) {
	var record = Record{}
	if err := snapshot.Unmarshal(&record); err != nil {
		return nil, err
	}
	if likes, ok := record["likes"]; ok {
		if likesMap, ok := likes.(map[string]interface{}); ok {
			record["likes_count"] = len(likesMap)
		}
	}
	return record, nil
}

func (s *service) UpdateLikes(ctx context.Context, resource Resource) error {
	databaseURL, err := resource.DatabaseURL()
	if err != nil {
		return err
	}
	refPath, err := resource.RefPath()
	if err != nil {
		return err
	}
	return RunTransaction(ctx, databaseURL, refPath, s.update)
}

func New() Service {
	return &service{}
}
