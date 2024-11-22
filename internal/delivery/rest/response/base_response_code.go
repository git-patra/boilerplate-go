package response

import "net/http"

type Code struct {
	HttpStatusCode int
	ErrorCode      string
	Message        string
}

var Ok Code = Code{http.StatusOK, "", "Request performed successfully"}
var Created Code = Code{http.StatusCreated, "", "Resource created successfully"}

var ValidationError Code = Code{http.StatusBadRequest, "VALIDATION_ERROR", "Request validation error"}
var GenericResourceNotFound Code = Code{http.StatusNotFound, "RESOURCE_NOT_FOUND", "Requested resource not found"}
var GenericServerError Code = Code{http.StatusInternalServerError, "SERVER_ERROR", "There is a problem processing your request. Error has been logged."}
