package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"

	. "tauterlyc/jsonapi/models"
)

type DocumentReader interface {
	json.Unmarshaler

	// Reads the resource from the data property and stores the result in the value pointed to by r, returns an error if the data property is a slice.
	Resource(r Resource) error
	// Reads the resources from the data property and stores the result in the value pointed to by r, returns an error if the data property is not a slice.
	Resources(r *[]Resource) error
}

type documentReader struct {
	Jsonapi  *JsonApi   `json:"jsonapi,omitempty"`
	Meta     Meta       `json:"meta,omitempty"`
	Links    Links      `json:"links,omitempty"`
	Included []Resource `json:"included,omitempty"`
	Errors   []TError   `json:"errors,omitempty"`

	data      Resource
	dataSlice []Resource
}

func Read(r *http.Request) (DocumentReader, error) {

	d := documentReader{}

	body := r.Body
	defer body.Close()

	src, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return &d, json.Unmarshal(src, &d)
}

func (d *documentReader) UnmarshalJSON(src []byte) error {

	base := struct {
		Jsonapi  *JsonApi   `json:"jsonapi,omitempty"`
		Links    Links      `json:"links,omitempty"`
		Included []Resource `json:"included,omitempty"`
		Errors   []TError   `json:"errors,omitempty"`
		Meta     Meta       `json:"meta,omitempty"`
	}{}

	if err := json.Unmarshal(src, &base); err != nil {
		return err
	}

	d.Jsonapi = base.Jsonapi
	d.Links = base.Links
	d.Errors = base.Errors
	d.Meta = base.Meta

	d.Included = make([]Resource, len(base.Included))
	copy(d.Included, base.Included)

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(src, &rawMap); err != nil {
		return err
	}

	data := rawMap["data"]

	if data == nil {
		return nil
	}

	if regexp.MustCompile(`^[\t\s\n]*{`).Match(data) {
		resource := NewResource("", "")
		if err := json.Unmarshal(data, &resource); err != nil {
			return err
		}
		d.data = resource
		return nil
	} else if regexp.MustCompile(`^[\t\s\n]*\[`).Match(data) {
		raw := make([]json.RawMessage, 0)
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}

		resources := make([]Resource, len(raw))

		for i, str := range raw {
			resource := NewResource("", "")
			if err := json.Unmarshal(str, &resource); err != nil {
				return err
			}
			resources[i] = resource
		}

		d.dataSlice = resources
		return nil
	}

	return errors.New("unable to read data")

}

func (d *documentReader) Resource(r Resource) error {
	// 1. Get the reflection Value of the argument
	rv := reflect.ValueOf(r)

	// 2. Validate it is a pointer and not nil
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return Error(ValueNotPointer, rv.Kind())
	}

	if d.dataSlice != nil {
		// Note: You can't "nil" the caller's pointer easily if it's a concrete struct
		return Error(DataNotResource)
	}

	// 3. Get the "Value" that the pointer points to
	target := rv.Elem()

	// 4. Get the source value
	source := reflect.ValueOf(d.data).Elem()

	// 5. Overwrite the memory at r with the data from d.data
	if !source.Type().AssignableTo(target.Type()) {
		return fmt.Errorf("type mismatch: cannot set %T to %T", d.data, r)
	}

	target.Set(source)

	return nil

}

func (d *documentReader) Resources(r *[]Resource) error {

	if d.data != nil {
		r = nil
		return Error(DataNotCollection)
	}

	for i, item := range d.dataSlice {

		x, ok := item.(Resource)
		*r = append(*r, x)

		if !ok {
			*r = nil
			return fmt.Errorf("unable to parse item at data[%v]", i)
		}

	}

	return nil

}
