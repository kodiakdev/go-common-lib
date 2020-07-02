package commonlib

import "go.mongodb.org/mongo-driver/mongo"

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
