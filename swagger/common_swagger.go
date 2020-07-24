package commonswagger

import (
	"fmt"

	commondto "github.com/kodiakdev/go-common-lib/dto"
)

var MsgDBProblem = GenerateSwaggerError(commondto.DatabaseErrorCode, commondto.DatabaseErrorExplanation)
var MsgUnknownProblem = GenerateSwaggerError(commondto.UnknownErrorCode, commondto.UnknownErrorExplanation)

//GenerateSwaggerError generate a string in a format of %d - %s
//use it for creating sample error message in swagger page
func GenerateSwaggerError(code int, explanation string) string {
	return fmt.Sprintf("%d - %s", code, explanation)
}

//CombineSwaggerError combine multiple error messages in one string
//use it for creating sample error message in swagger page
func CombineSwaggerError(messages ...string) string {
	finalMsg := ""
	for i := 0; i < len(messages); i++ {
		if len(messages)-1 == i {
			finalMsg += messages[i]
		} else {
			finalMsg += messages[i] + ` | `
		}
	}
	return finalMsg
}
