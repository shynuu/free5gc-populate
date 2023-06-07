package runtime

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// This function come from https://github.com/free5gc/MongoDBLibrary/blob/main/api_mongoDB.go (License Apache 2)
// with new parameters "client" and "dbName", and a change on the return type
func RestfulAPIPost(client *mongo.Client, dbName string, collName string, filter bson.M, postData map[string]interface{}) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("mongo client is nil")
	}
	collection := client.Database(dbName).Collection(collName)

	var checkItem map[string]interface{}
	collection.FindOne(context.TODO(), filter).Decode(&checkItem)

	if checkItem == nil {
		collection.InsertOne(context.TODO(), postData)
		return false, nil
	} else {
		collection.UpdateOne(context.TODO(), filter, bson.M{"$set": postData})
		return true, nil
	}
}

// This function come from https://github.com/free5gc/MongoDBLibrary/blob/main/api_mongoDB.go (License Apache 2)
// with new parameters "client" and "dbName", and a change on the return type
func RestfulAPIPostMany(client *mongo.Client, dbName string, collName string, filter bson.M, postDataArray []interface{}) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("mongo client is nil")
	}
	collection := client.Database(dbName).Collection(collName)

	collection.InsertMany(context.TODO(), postDataArray)
	return false, nil
}
