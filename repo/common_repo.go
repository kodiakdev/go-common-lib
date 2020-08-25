package commonrepo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IDBOperation interface {
	InsertOne(data interface{}) (*mongo.InsertOneResult, error)
	InsertMany(data []interface{}) (*mongo.InsertManyResult, error)
	FindOne(filter, impl interface{}) (interface{}, error)
	FindOneAndUpdate(filter, update interface{}) (*mongo.UpdateResult, error)
	Find(filter, impl interface{}) (interface{}, error)
	FindPagedSorted(pagingSortingReq PagingSortingRequest, filter, impl interface{}) (interface{}, *mongopagination.PaginationData, error)
	Count(filter interface{}) (int64, error)
	IsExist(filter interface{}) (bool, error)
	InsertOneAtColl(collection string, data interface{}) (*mongo.InsertOneResult, error)
	InsertManyAtColl(collection string, data []interface{}) (*mongo.InsertManyResult, error)
	FindOneAtColl(collection string, filter, impl interface{}) (interface{}, error)
	FindOneAndUpdateAtColl(collection string, filter, update interface{}) (*mongo.UpdateResult, error)
	FindAtColl(collection string, filter, impl interface{}) (interface{}, error)
	FindAtCollPagedSorted(collection string, pagingSortingReq PagingSortingRequest, filter, impl interface{}) (interface{}, *mongopagination.PaginationData, error)
	CountAtColl(collection string, filter interface{}) (int64, error)
	IsExistAtColl(collection string, filter interface{}) (bool, error)
}

//PagingSortingRequest specify paging and sorting mechanism
//SortByField fill with name of the field to be taken by sorter
//IsSortDesc leave empty (default false) if sort in ascending order
//LimitPerPage number of items per page
//Page page index (start at 1)
type PagingSortingRequest struct {
	SortByField  string
	IsSortDesc   bool
	LimitPerPage int64
	Page         int64
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

func (dbOp *DBOperation) InsertMany(data []interface{}) (*mongo.InsertManyResult, error) {
	return dbOp.InsertManyAtColl(dbOp.defaultCollection, data)
}

func (dbOp *DBOperation) FindOne(filter, impl interface{}) (interface{}, error) {
	return dbOp.FindOneAtColl(dbOp.defaultCollection, filter, impl)
}

func (dbOp *DBOperation) FindOneAndUpdate(filter, update interface{}) (*mongo.UpdateResult, error) {
	return dbOp.FindOneAndUpdateAtColl(dbOp.defaultCollection, filter, update)
}

func (dbOp *DBOperation) Find(filter, impl interface{}) (interface{}, error) {
	return dbOp.FindAtColl(dbOp.defaultCollection, filter, impl)
}

func (dbOp *DBOperation) FindPagedSorted(pagingSortingReq PagingSortingRequest, filter, impl interface{}) (interface{}, *mongopagination.PaginationData, error) {
	return dbOp.FindAtCollPagedSorted(dbOp.defaultCollection, pagingSortingReq, filter, impl)
}

func (dbOp *DBOperation) Count(filter interface{}) (int64, error) {
	return dbOp.CountAtColl(dbOp.defaultCollection, filter)
}

func (dbOp *DBOperation) IsExist(filter interface{}) (bool, error) {
	return dbOp.IsExistAtColl(dbOp.defaultCollection, filter)
}

func (dbOp *DBOperation) InsertOneAtColl(collection string, data interface{}) (*mongo.InsertOneResult, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()

	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(collection)
	inserted, err := coll.InsertOne(ctx, data)
	return inserted, err
}

func (dbOp *DBOperation) InsertManyAtColl(collection string, data []interface{}) (*mongo.InsertManyResult, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()

	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(collection)
	inserted, err := coll.InsertMany(ctx, data)
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

func (dbOp *DBOperation) FindOneAndUpdateAtColl(collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()

	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(collection)
	res, err := coll.UpdateOne(ctx, filter, update)
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

//FindAtCollPagedSorted find at specified collection and return as paged and sorted result
func (dbOp *DBOperation) FindAtCollPagedSorted(collection string, pagingSortingReq PagingSortingRequest, filter, impl interface{}) (interface{}, *mongopagination.PaginationData, error) {
	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(collection)
	paginatedData, err := mongopagination.New(coll).
		Limit(pagingSortingReq.LimitPerPage).
		Page(pagingSortingReq.Page).
		Sort(pagingSortingReq.SortByField, dbOp.convertToSortType(pagingSortingReq.IsSortDesc)).
		Filter(filter).
		Find()

	if err != nil {
		logrus.Errorf("Failed to perform query with sort at coll %s. PagingSortingReq: %v | Error: %v", collection, pagingSortingReq, err)
		return nil, nil, err
	}

	if paginatedData.Data == nil {
		return nil, nil, nil
	}

	err = dbOp.decodeBsonRaws(paginatedData.Data, impl)
	if err != nil {
		logrus.Errorf("Failed to perform dbOp.decodeBsonRaws(paginatedData.Data, impl). Impl: %v | Error: %v", impl, err)
		return nil, nil, err
	}

	return impl, &paginatedData.Pagination, nil
}

func (dbOp *DBOperation) CountAtColl(collection string, filter interface{}) (int64, error) {
	ctx, cf := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cf()
	coll := dbOp.dbClient.Database(dbOp.databaseName).Collection(collection)
	counter, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return counter, nil
}
func (dbOp *DBOperation) IsExistAtColl(collection string, filter interface{}) (bool, error) {
	exist, err := dbOp.CountAtColl(collection, filter)
	if err != nil {
		return false, err
	}
	return (exist > 0), nil
}

//decodeBsonRaws iterate the bson.Raw array and decode to impl
func (dbOp *DBOperation) decodeBsonRaws(rawData []bson.Raw, impl interface{}) error {
	resultsVal := reflect.ValueOf(impl)
	if resultsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer %v", resultsVal.Type())
	}
	// get the value that the pointer v points to.
	sliceVal := resultsVal.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return errors.New("can't fill non-slice value")
	}
	elemType := sliceVal.Type().Elem()
	var index int
	for _, w := range rawData {
		if sliceVal.Len() == index {
			// slice is full
			newElem := reflect.New(elemType)
			sliceVal = reflect.Append(sliceVal, newElem.Elem())
			sliceVal = sliceVal.Slice(0, sliceVal.Cap())
		}

		currElem := sliceVal.Index(index).Addr().Interface()
		err := bson.Unmarshal(w, currElem)
		if err != nil {
			return fmt.Errorf("can't unmarshal. Err: %v", err)
		}
		index++
	}
	resultsVal.Elem().Set(sliceVal.Slice(0, index))
	return nil
}

func (dbOp *DBOperation) convertToSortType(sortDesc bool) int {
	if sortDesc {
		return -1
	}
	return 1
}
