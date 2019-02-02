package cloud_function

import (
	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

const (
	ScopeDb    = "https://www.googleapis.com/auth/firebase.database"
	ScopeEmail = "https://www.googleapis.com/auth/userinfo.email"
)

//NewApp returns new app
func NewApp(ctx context.Context, URL string) (*firebase.App, error) {
	options := []option.ClientOption{
		option.WithScopes(ScopeEmail, ScopeDb),
	}
	return firebase.NewApp(ctx, &firebase.Config{
		DatabaseURL: URL,
	}, options...)

}

//NewDb returns new db
func NewDb(ctx context.Context, URL string) (*db.Client, error) {
	app, err := NewApp(ctx, URL)
	if err != nil {
		return nil, err
	}
	return app.Database(ctx)
}
