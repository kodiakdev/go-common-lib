package commonrepo

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type IDBOperation interface {
	InsertOne(tableName string, data interface{}) (*mongo.InsertOneResult, error)
	FindOne(tableName string, filter, impl interface{}) (interface{}, error)
	FindOneAndUpdate(tableName string, filter, data interface{}) (*mongo.UpdateResult, error)
	Find(tableName string, filter, impl interface{}) (interface{}, error)
}

type DBOperation struct {
	dbClient     *mongo.Client
	databaseName string
}

func NewDBOperation(dbClient *mongo.Client, databaseName string) *DBOperation {
	return &DBOperation{
		dbClient:     dbClient,
		databaseName: databaseName,
	}
}

func (dbOp *DBOperation) InsertOne(tableName string, data interface{}) (*mongo.InsertOneResult, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()

	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(tableName)
	inserted, err := coll.InsertOne(ctx, data)
	return inserted, err
}

func (dbOp *DBOperation) FindOne(tableName string, filter, impl interface{}) (interface{}, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()

	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(tableName)
	res := coll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	err := res.Decode(impl)
	if err != nil {
		logrus.Errorf("Failed to perform res.Decode(&impl). Impl: %v | Error: %v", impl, err)
		return nil, err
	}

	return impl, nil
}

func (dbOp *DBOperation) FindOneAndUpdate(tableName string, filter, data interface{}) (*mongo.UpdateResult, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()

	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(tableName)
	res, err := coll.UpdateOne(ctx, filter, bson.D{{"$set", data}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return res, err
	}

	return res, nil
}

func (dbOp *DBOperation) Find(tableName string, filter, impl interface{}) (interface{}, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()
	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(tableName)
	curr, err := coll.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	defer curr.Close(ctx)

	err = curr.All(ctx, impl)
	if err != nil {
		logrus.Errorf("Failed to perform curr.All(ctx, impl). Impl: %v | Error: %v", impl, err)
		return nil, err
	}

	return impl, nil
}
