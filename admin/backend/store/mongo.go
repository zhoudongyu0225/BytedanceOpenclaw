package store

import (
	"admin-backend/config"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func InitMongo(cfg config.MongoConfig) error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return err
	}
	// 测试连接
	if err := client.Ping(context.TODO(), nil); err != nil {
		return err
	}
	MongoClient = client
	MongoDB = client.Database(cfg.Database)
	return nil
}

func CloseMongo() {
	if MongoClient != nil {
		MongoClient.Disconnect(context.TODO())
	}
}

func GetCollection(name string) *mongo.Collection {
	return MongoDB.Collection(name)
}
