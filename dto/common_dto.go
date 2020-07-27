package commondto

//error code mapping:
// aabbccc
// aa: service/component identity (10 reserved for general error)
// bb: error identity (99: general error, 00: bad request, 01: unauthorized, 03: forbidden, 04: not found, 09: conflict)
// ccc: detailed error identity
const (
	UnknownErrorCode                  = 1099001
	UnknownErrorExplanation           = "Unknown error occured. Contact the technical support."
	FailedParseRequestBodyCode        = 1099002
	FailedParseRequestBodyExplanation = "Failed to parse request body. Make sure your request conform with api documentation."
	DatabaseErrorCode                 = 1099003
	DatabaseErrorExplanation          = "Error occured during read/write database"
	FeatureNotAvailableCode           = 1099004
	FeatureNotAvailableExplanation    = "Feature not available yet"
	IncompleteInputCode               = 1099005
	IncompleteInputExplanation        = "At least one mandatory input not provided. Check API documentation!"
)
