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

//NewClient returns new firestore client
func NewClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	options := []option.ClientOption{
		option.WithScopes(ScopeEmail, DatastoreDb, ScopeDb),
	}

	return firestore.NewClient(ctx, projectID, options...)
}
