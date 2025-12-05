package infra

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const maxPoolSize = 500

func InitMongoDBClient(uri string) (*mongo.Client, error) {
	monitor := &event.PoolMonitor{
		Event: HandlePoolMonitor,
	}
	opts := options.Client().ApplyURI(uri).
		SetMaxPoolSize(maxPoolSize).
		SetPoolMonitor(monitor)
	return mongo.Connect(context.Background(), opts)
}

//nolint:gocritic
func HandlePoolMonitor(evt *event.PoolEvent) {
	switch evt.Type {
	case event.PoolClosedEvent:
		fmt.Println("DB connection closed.")
	}
}
