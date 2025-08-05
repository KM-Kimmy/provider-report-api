package clienterrors

import "errors"

var ErrAlreadyDeleted = errors.New("object is already deleted")
var ErrInsurerDeleted = errors.New("insurer is deleted or does not exist")
var ErrProviderDeleted = errors.New("provier is deleted or does not exist")
var ErrNotFound = errors.New("object is deleted or does not exist")
var ErrDataNotFound = errors.New("data not found")

var Errn = errors.New("company is deleted or does not exist")

type InternalserverErrorResponse struct {
	Message string `json:"message" example:"Internal server error, cannot process the request."`
}

type BadrequestErrorResponse struct {
	Message string `json:"message"  example:"Invalid input parameters."`
}

type ConflictErrorResponse struct {
	Message string `json:"message" example:"object is already deleted"`
}
