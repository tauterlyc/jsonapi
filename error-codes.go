package jsonapi

import "net/http"

type ErrorType struct {
	code    int
	status  int
	title   string
	details string
}

var FormReadError = ErrorType{
	code:    100,
	status:  http.StatusBadRequest,
	title:   "form data unreadable",
	details: "%v",
}

var MultipartFormReadError = ErrorType{
	code:    102,
	status:  http.StatusBadRequest,
	title:   "multipart form data unreadable",
	details: "%v",
}

var JsonReadError = ErrorType{
	code:    103,
	status:  http.StatusBadRequest,
	title:   "json data unreadable",
	details: "%v",
}

var ForbiddenRelationship = ErrorType{
	code:    104,
	status:  http.StatusForbidden,
	title:   "invalid relationship",
	details: "resource relationship type of \"%s\" not allowed",
}

var ForbiddenToManyRelationship = ErrorType{
	code:    105,
	status:  http.StatusForbidden,
	title:   "to many relationship unsupported",
	details: "resource relationship type of \"%s\" cannot be a to-many relationship",
}

var ReferenceNotExist = ErrorType{
	code:    106,
	status:  http.StatusNotFound,
	title:   "related resource does not exist",
	details: "unable to find resource with ID \"%s\" and type \"%s\" for relationship \"%s\"",
}

var ResourceNotExist = ErrorType{
	code:    107,
	status:  http.StatusNotFound,
	title:   "resource does not exist",
	details: "resource with type \"%s\" and ID \"%s\" does not exist",
}

var RelationshipNotExist = ErrorType{
	code:    108,
	status:  http.StatusNotFound,
	title:   "resource relationship unsupported",
	details: "\"%s\" do not support \"%s\" relationships",
}

var JSONRequestReadError = ErrorType{
	code:    109,
	status:  http.StatusBadRequest,
	title:   "cannot read json data",
	details: "%v",
}

var ResourceAlreadyExists = ErrorType{
	code:    110,
	status:  http.StatusConflict,
	title:   "resource already exists",
	details: "resource with type \"%s\" and id of \"%s\" already exists",
}

var DatabaseReadFailed = ErrorType{
	code:    111,
	status:  http.StatusInternalServerError,
	title:   "database read failed",
	details: "unable to read %s from database: %w",
}

var DatabaseWriteFailed = ErrorType{
	code:    112,
	status:  http.StatusInternalServerError,
	title:   "database write failed",
	details: "unable to write %s to database: %w",
}

var DataNotCollection = ErrorType{
	code:    113,
	status:  http.StatusUnprocessableEntity,
	title:   "data not array",
	details: "data must be an array of resource objects",
}

var DataNotResource = ErrorType{
	code:   114,
	status: http.StatusUnprocessableEntity,
	title:  "data must be a resource object",
}

var UnsupportedQuery = ErrorType{
	code:    115,
	status:  http.StatusBadRequest,
	title:   "query unsupported",
	details: "unsupported query: %s",
}

var MissingGetter = ErrorType{
	code:    116,
	status:  http.StatusInternalServerError,
	title:   "jsonapi: missing include getter",
	details: "jsonapi: missing include getter for \"%s\"",
}

var ValueNotPointer = ErrorType{
	code:    117,
	status:  http.StatusInternalServerError,
	title:   "jsonapi: value must be a pointer",
	details: "jsonapi: value must be a pointer, got \"%v\"",
}
