package jsonapi

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/tauterlyc/jsonapi/models"
)

func ResourceFromRequest(r *http.Request) (models.Resource, error) {

	upsert, err := Read(r)
	if err != nil {
		return nil, err
	}

	resource := models.NewResource("", "")
	err = upsert.Resource(resource)
	if err != nil {
		return nil, errors.Join(errors.New("failed to read resource from the request body"), err)
	}

	return resource, nil

}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	New(r).AddError(errors.New("method now allowed")).WriteHeader(http.StatusMethodNotAllowed).Write(w)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	New(r).AddError(errors.New("not implemented")).WriteHeader(http.StatusNotImplemented).Write(w)
}

func ReadQueryParams(query url.Values) (map[string][]string, error) {
	var fields map[string][]string = nil

	for key, values := range query {
		if strings.HasPrefix(key, "fields[") && strings.HasSuffix(key, "]") {
			if fields == nil {
				fields = map[string][]string{}
			}

			innerKey := key[7 : len(key)-1] // Extracts "TYPE" from "fields[TYPE]"

			fields[innerKey] = strings.Split(values[0], ",")
		}
	}

	return fields, nil
}
