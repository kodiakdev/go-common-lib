package commonresp

import (
	"net/http"

	commondto "github.com/kodiakdev/go-common-lib/dto"

	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
)

// RequestResponse holds all value of request and response
type RequestResponse struct {
	Req          *restful.Request
	Resp         *restful.Response
	HTTPStatus   int
	ResponseBody interface{}
}

//ServiceErrorResponse response commonerr for non 2xx
type ServiceErrorResponse struct {
	Code        int                         `json:"code"`
	Explanation string                      `json:"explanation"`
	Causes      []ServiceErrorCauseResponse `json:"causes,omitempty"`
}

//ServiceErrorCauseResponse struct explaining error cause
type ServiceErrorCauseResponse struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

//Write perform write the response as JSON
func write(comm *RequestResponse) {
	err := comm.Resp.WriteHeaderAndJson(
		comm.HTTPStatus,
		comm.ResponseBody,
		restful.MIME_JSON,
	)
	if err != nil {
		logrus.Warnf("Unable to write response. Error was %v", err)
	}
}

func RespondRequestParsingFail(err error, req *restful.Request, resp *restful.Response) {
	logrus.Warnf("Failed to read entity. Error: %v", err)
	errorResponseBody := ServiceErrorResponse{
		Code:        commondto.FailedParseRequestBodyCode,
		Explanation: commondto.FailedParseRequestBodyExplanation,
	}
	Respond(errorResponseBody, http.StatusBadRequest, req, resp)
}

func RespondDatabaseError(err error, req *restful.Request, resp *restful.Response) {
	errorResponseBody := ServiceErrorResponse{
		Code:        commondto.DatabaseErrorCode,
		Explanation: commondto.DatabaseErrorExplanation,
	}
	Respond(errorResponseBody, http.StatusInternalServerError, req, resp)
}

func RespondIncompleteInput(err error, req *restful.Request, resp *restful.Response, causes ...ServiceErrorCauseResponse) {
	errorResponseBody := ServiceErrorResponse{
		Code:        commondto.IncompleteInputCode,
		Explanation: commondto.IncompleteInputExplanation,
		Causes:      causes,
	}
	Respond(errorResponseBody, http.StatusBadRequest, req, resp)
}

func RespondUnknownError(err error, req *restful.Request, resp *restful.Response) {
	logrus.Errorf("Error occured with unknown reason. Error: %v", err)
	errorResponseBody := ServiceErrorResponse{
		Code:        commondto.UnknownErrorCode,
		Explanation: commondto.UnknownErrorExplanation,
	}
	Respond(errorResponseBody, http.StatusInternalServerError, req, resp)
}

func Respond(body interface{}, httpStatus int, req *restful.Request, resp *restful.Response) {
	write(&RequestResponse{
		Req:          req,
		Resp:         resp,
		HTTPStatus:   httpStatus,
		ResponseBody: body,
	})
}
