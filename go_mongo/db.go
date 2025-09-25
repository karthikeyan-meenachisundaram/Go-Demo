package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoURI returns the URI string (hidden from main.go)
func MongoURI() string {
	// You can also read this from an env variable here if you want
	return "mongodb+srv://Karthikeyan:Hema%401199@mycluster.5oolqvy.mongodb.net/"
}

// Connect returns Mongo client, context, and cancel func.
// Caller must defer cancel() and client.Disconnect(ctx).
func Connect(timeoutSeconds int) (*mongo.Client, context.Context, context.CancelFunc, error) {
	uri := MongoURI() // get URI from this file

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}

	// optional: ping
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		cancel()
		return nil, nil, nil, err
	}

	fmt.Println("✅ Connected to MongoDB")
	return client, ctx, cancel, nil
}
