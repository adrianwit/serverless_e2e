package fb

import (
	"context"
	"firebase.google.com/go/db"
)

//RunTransaction runs transactions
func RunTransaction(ctx context.Context, URL, refPath string, updateFunc func(nodeSnapshot db.TransactionNode) (interface{}, error)) error {
	dbc, err := NewDb(ctx, URL)
	if err != nil {
		return err
	}
	ref := dbc.NewRef(refPath)
	return ref.Transaction(ctx, updateFunc)
}
