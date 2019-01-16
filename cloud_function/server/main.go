package main

import (
	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/viant/toolbox"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"os"
)

const (
	prodAddr = "https://abstractdb-154a9.firebaseio.com"
	// Scope is the Oauth2 scope for the service.
	ScopeDb    = "https://www.googleapis.com/auth/firebase.database"
	ScopeEmail = "https://www.googleapis.com/auth/userinfo.email"
)

func main() {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/Users/awitas/.secret/am.json")

	ctx := context.Background()
	options := []option.ClientOption{
		option.WithScopes(ScopeEmail, ScopeDb),
	}
	app, err := firebase.NewApp(ctx, &firebase.Config{
		DatabaseURL: prodAddr,
	}, options...)
	if err != nil {
		log.Fatal(err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatal(err)
	}
	ref := client.NewRef("users/100/likes")
	err = ref.Transaction(ctx, func(currentSnapshot db.TransactionNode) (result interface{}, err error) {
		var counter = 0
		var previous interface{}
		if err = currentSnapshot.Unmarshal(&previous); err != nil {
			return nil, err
		}
		if previous != nil {
			if intValue, err := toolbox.ToInt(previous); err == nil {
				counter = intValue
			}
		}
		counter++
		return counter, nil

	})

	if err != nil {
		log.Fatal(err)
	}

}
