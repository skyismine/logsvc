package storage

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type StorageMongodb struct {
	mongoclient *mongo.Client
}

func NewStorageMongodb(addr string) *StorageMongodb {
	mongodb := new(StorageMongodb)
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongodb.mongoclient, err = mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", addr)))
	if err != nil {
		logs.Error("init mongo.Connect error", err)
		return nil
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = mongodb.mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		logs.Error("init mongoclient.Ping error", err)
		return nil
	}
	return mongodb
}

func (store *StorageMongodb) Save(msg *Logmsg) error {
	collection := store.mongoclient.Database(msg.App).Collection("logs")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func (store *StorageMongodb) Close() {
	_ = store.mongoclient.Disconnect(nil)
}