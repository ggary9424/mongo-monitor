package mongowrapper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateClient(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(uri),
		nil,
	)

	return client, err
}
