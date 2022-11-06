package mongodb

import (
	"context"
	"fmt"
	"go-cloud-camp/internal/config"
	"go-cloud-camp/internal/logging"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongBackend struct
type MongoBackend struct {
	client *mongo.Client
	mdb    *mongo.Database
	logger *logging.Logger
}

func Create(cfg *config.StorageParams, logger *logging.Logger) (*MongoBackend, error) {
	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:%d/?maxPoolSize=%d&w=majority",
		cfg.MongoDB.User,
		cfg.MongoDB.Pass,
		cfg.MongoDB.Host,
		cfg.MongoDB.Port,
		cfg.MongoDB.MaxPoolSize,
	)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}

	// Ping mongodb server
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}

	logger.Info("connected to mongodb backend")

	return &MongoBackend{
		client: client,
		mdb:    client.Database(cfg.MongoDB.Database),
		logger: logger,
	}, nil
}

// Close function
func (mb *MongoBackend) Close(ctx context.Context) error {
	return mb.client.Disconnect(ctx)
}
