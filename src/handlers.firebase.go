package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"fmt"
	"google.golang.org/api/option"
)

var dataBaseApp *firebase.App
var ctx context.Context

func initializeApp(credPath []byte, databaseURL string) error {
	opt := option.WithCredentialsJSON(credPath)
	//option.WithCredentialsJSON(?)
	ctx = context.Background()
	conf := &firebase.Config{
		DatabaseURL: databaseURL,
	}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	dataBaseApp = app
	return err
}

func initDB() (*db.Client, error) {
	client, err := dataBaseApp.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing database client: %v", err)
	}

	return client, nil
}

func readEntry(client *db.Client, path string) (map[string]interface{}, error) {
	ref := client.NewRef(path)
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		return nil, fmt.Errorf("error reading from database: %v", err)
	}
	return data, nil
}

func readValue(client *db.Client, path string) (interface{}, error) {
	ref := client.NewRef(path)
	var data interface{}
	if err := ref.Get(ctx, &data); err != nil {
		return nil, fmt.Errorf("error reading from database: %v", err)
	}
	return data, nil
}

func pathExists(client *db.Client, path string) bool {
	ref := client.NewRef(path)
	var data map[string]interface{}
	err := ref.GetShallow(ctx, &data)
	return err == nil && len(data) != 0
}

func setEntry(client *db.Client, path string, value interface{}) error {
	ref := client.NewRef(path)
	if err := ref.Set(ctx, value); err != nil {
		return fmt.Errorf("error setting to database: %v", err)
	}
	return nil
}

func deletePath(client *db.Client, path string) error {
	ref := client.NewRef(path)
	if err := ref.Set(ctx, nil); err != nil {
		return fmt.Errorf("error setting to database: %v", err)
	}
	return nil
}

func pushEntry(client *db.Client, path string, value interface{}) (string, error) {
	ref := client.NewRef(path)
	loc, err := ref.Push(ctx, value)
	if err != nil {
		return "", fmt.Errorf("error setting to database: %v", err)
	}
	return loc.Key, nil
}
