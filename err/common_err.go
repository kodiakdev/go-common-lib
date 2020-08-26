package commonerr

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrDatabaseProblem = errors.New("dbproblem")
var ErrUnknownProblem = errors.New("unknown")
var ErrUniqueConstraintViolation = errors.New("uniqueconstraint")
var ErrRecordNotFound = errors.New("recnotfound")
var ErrNotImplementedYet = errors.New("notimplementedyet")
// ErrInvalidRequest should be thrown when input from api is not valid
var ErrInvalidRequest = errors.New("invalidreq")

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
