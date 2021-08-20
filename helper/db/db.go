package db

import (
	"AuthtggO/logHelper"
	"AuthtggO/utils"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DbClient struct {
	Mongo *mongo.Client
}

var DatabaseClient *DbClient

func initClient() *mongo.Client {
	logger := logHelper.GetLogger()
	logger.Debugf("[DB] Connection: %s\n", utils.GetDbUri())
	client, err := mongo.NewClient(options.Client().ApplyURI(utils.GetDbUri()))
	if err != nil {
		logger.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	return client
}

func InitDbClient() {
	DatabaseClient = &DbClient{Mongo: initClient()}
}

