package mongowrapper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateClient(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
		nil,
	)

	return client, err
}
