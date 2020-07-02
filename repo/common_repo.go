package commonrepo

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type IDBOperation interface {
	InsertOne(data interface{}) (*mongo.InsertOneResult, error)
	FindOne(filter, impl interface{}) (interface{}, error)
	FindOneAndUpdate(filter, data interface{}) (*mongo.UpdateResult, error)
	Find(filter, impl interface{}) (interface{}, error)
	InsertOneAtColl(collection string, data interface{}) (*mongo.InsertOneResult, error)
	FindOneAtColl(collection string, filter, impl interface{}) (interface{}, error)
	FindOneAndUpdateAtColl(collection string, filter, data interface{}) (*mongo.UpdateResult, error)
	FindAtColl(collection string, filter, impl interface{}) (interface{}, error)
}

type DBOperation struct {
	dbClient          *mongo.Client
	databaseName      string
	defaultCollection string
}

func NewDBOperation(dbClient *mongo.Client, databaseName, defaultCollection string) *DBOperation {
	return &DBOperation{
		dbClient:          dbClient,
		databaseName:      databaseName,
		defaultCollection: defaultCollection,
	}
}

func (dbOp *DBOperation) InsertOne(data interface{}) (*mongo.InsertOneResult, error) {
	return dbOp.InsertOneAtColl(dbOp.defaultCollection, data)
}

func (dbOp *DBOperation) FindOne(filter, impl interface{}) (interface{}, error) {
	return dbOp.FindOneAtColl(dbOp.defaultCollection, filter, impl)
}

func (dbOp *DBOperation) FindOneAndUpdate(filter, data interface{}) (*mongo.UpdateResult, error) {
	return dbOp.FindOneAndUpdateAtColl(dbOp.defaultCollection, filter, data)
}

func (dbOp *DBOperation) Find(filter, impl interface{}) (interface{}, error) {
	return dbOp.FindAtColl(dbOp.defaultCollection, filter, impl)
}

func (dbOp *DBOperation) InsertOneAtColl(collection string, data interface{}) (*mongo.InsertOneResult, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()

	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(collection)
	inserted, err := coll.InsertOne(ctx, data)
	return inserted, err
}

func (dbOp *DBOperation) FindOneAtColl(collection string, filter, impl interface{}) (interface{}, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()

	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(collection)
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

func (dbOp *DBOperation) FindOneAndUpdateAtColl(collection string, filter, data interface{}) (*mongo.UpdateResult, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()

	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(collection)
	res, err := coll.UpdateOne(ctx, filter, bson.D{{"$set", data}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return res, err
	}

	return res, nil
}

func (dbOp *DBOperation) FindAtColl(collection string, filter, impl interface{}) (interface{}, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()
	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(collection)
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
