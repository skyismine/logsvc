package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var mongoclient *mongo.Client

func init() {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoclient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalln("init mongo.Connect error", err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln("init mongoclient.Ping error", err)
	}
}

func mgodbinsert(app string, msg interface{}) {
	collection := mongoclient.Database(app).Collection("logs")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, msg)
	if err != nil {
		log.Println("mgodbinsert collection.InsertOne error", err)
		return
	}
}
