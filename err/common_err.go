package commonerr

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrDatabaseProblem = errors.New("dbproblem")
var ErrUnknownProblem = errors.New("unknown")

const (
	UnknownErrorCode                  = 1000001
	UnknownErrorExplanation           = "Unknown error occured. Contact the technical support."
	FailedParseRequestBodyCode        = 1000002
	FailedParseRequestBodyExplanation = "Failed to parse request body. Make sure your request conform with api documentation."
	DatabaseErrorCode                 = 1000003
	DatabaseErrorExplanation          = "Error occured during read/write database"
)

//IsDuplicateErr check if err is a mongodb duplication error
//duplication error happens in column with unique constraint
func IsDuplicateErr(err error) bool {
	if errWriteException, ok := err.(mongo.WriteException); ok {
		if len(errWriteException.WriteErrors) > 0 {
			writeErr := errWriteException.WriteErrors[0]
			if writeErr.Code == 11000 {
				return true
			}
		}
	}
	return false
}
