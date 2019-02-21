package cloud_function

import (
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

const (
	ScopeDb     = "https://www.googleapis.com/auth/firebase.database"
	DatastoreDb = "https://www.googleapis.com/auth/datastore"
	ScopeEmail  = "https://www.googleapis.com/auth/userinfo.email"
)

var client *firestore.Client

//NewClient returns new firestore client
func NewClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	if client != nil {
		return client, nil
	}
	options := []option.ClientOption{
		option.WithScopes(ScopeEmail, DatastoreDb, ScopeDb),
	}
	var err error
	client, err = firestore.NewClient(ctx, projectID, options...)
	return client, err
}
