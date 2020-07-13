package testutilRepo

import (
	mongopagination "github.com/gobeam/mongo-go-pagination"
	commonrepo "github.com/kodiakdev/go-common-lib/repo"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockDBOperation struct {
	mock.Mock
	defaultCollection string
}

func (dbOp *MockDBOperation) InsertOne(data interface{}) (*mongo.InsertOneResult, error) {
	return dbOp.InsertOneAtColl(dbOp.defaultCollection, data)
}

func (dbOp *MockDBOperation) FindOne(filter, impl interface{}) (interface{}, error) {
	return dbOp.FindOneAtColl(dbOp.defaultCollection, filter, impl)
}

func (dbOp *MockDBOperation) FindOneAndUpdate(filter, update interface{}) (*mongo.UpdateResult, error) {
	return dbOp.FindOneAndUpdateAtColl(dbOp.defaultCollection, filter, update)
}

func (dbOp *MockDBOperation) Find(filter, impl interface{}) (interface{}, error) {
	return dbOp.FindAtColl(dbOp.defaultCollection, filter, impl)
}

func (dbOp *MockDBOperation) FindPagedSorted(pagingSortingReq commonrepo.PagingSortingRequest, filter, impl interface{}) (interface{}, *mongopagination.PaginationData, error) {
	return dbOp.FindAtCollPagedSorted(dbOp.defaultCollection, pagingSortingReq, filter, impl)
}

func (dbOp *MockDBOperation) InsertOneAtColl(collection string, data interface{}) (*mongo.InsertOneResult, error) {
	args := dbOp.Called(collection, data)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (dbOp *MockDBOperation) FindOneAtColl(collection string, filter, impl interface{}) (interface{}, error) {
	args := dbOp.Called(collection, filter, impl)
	return args.Get(0), args.Error(1)
}

func (dbOp *MockDBOperation) FindOneAndUpdateAtColl(collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	args := dbOp.Called(collection, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (dbOp *MockDBOperation) FindAtColl(collection string, filter, impl interface{}) (interface{}, error) {
	args := dbOp.Called(collection, filter, impl)
	return args.Get(0), args.Error(1)
}

func (dbOp *MockDBOperation) FindAtCollPagedSorted(collection string, pagingSortingReq commonrepo.PagingSortingRequest, filter, impl interface{}) (interface{}, *mongopagination.PaginationData, error) {
	args := dbOp.Called(collection, pagingSortingReq, filter, impl)
	return args.Get(0), args.Get(1).(*mongopagination.PaginationData), args.Error(2)
}
