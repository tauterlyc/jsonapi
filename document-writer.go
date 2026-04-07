package jsonapi

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	. "github.com/tauterlyc/jsonapi/models"
)

type DocumentWriter interface {
	json.Marshaler

	// Adds a single resource to the data property, calling will overwrite any previous values.
	AddResource(Resource) DocumentWriter
	// Writes a resource slice to the data property, calling will overwrite any previous values.
	AddResources([]Resource) DocumentWriter

	// Manually include a Resource in the document.
	Include(Resource) DocumentWriter

	// Sets a flag internally so that when `Write` is called the document writer will loop through each
	// relationship in the data and include it in the document if it matches the query and the getter is
	// defined.
	IncludeRelated() DocumentWriter

	LinkTo(name string, value string) DocumentWriter
	SetMeta(key string, value interface{}) DocumentWriter

	// AddError adds one or more errors to the document writer.
	//
	// If the error is nil, the method returns the writer unchanged.
	// If an error is present, any existing data payloads (Data, Included)
	// are cleared to ensure the document reflects an error state.
	//
	// If the error implements the errorGroup interface, it is unwrapped and
	// all contained errors are appended individually. Otherwise, the error
	// is appended as-is.
	AddError(error) DocumentWriter

	// Sorts the data based on the request query
	Sortable() DocumentWriter
	// Filters attribute fields based on the request query
	FieldSet() DocumentWriter

	WriteHeader(int) DocumentWriter
	SetHeader(key string, value string) DocumentWriter

	// Write the document to the ResponseWriter
	Write(http.ResponseWriter)

	// Returns true if `AddError` has been called with a value that is not nil
	HasErrors() bool
}

type documentWriter struct {
	Jsonapi  *JsonApi   `json:"jsonapi,omitempty"`
	Meta     Meta       `json:"meta,omitempty"`
	Links    Links      `json:"links,omitempty"`
	Included []Resource `json:"included,omitempty"`

	data      Resource
	dataSlice []Resource
	errors    []error

	useSort           bool
	useSparseFieldSet bool
	useInclude        bool

	req         *http.Request
	writeCalled bool
	status      int
	headers     map[string]string
}

func New(r *http.Request) DocumentWriter {
	return &documentWriter{
		req:     r,
		Jsonapi: &JsonApi{Version: "1.1"}, // Todo: remove this.
	}
}

func (d *documentWriter) AddResource(r Resource) DocumentWriter {

	if d.errors != nil {
		return d
	}

	d.dataSlice = nil

	d.data = r

	return d
}

func (d *documentWriter) AddResources(r []Resource) DocumentWriter {

	if d.errors != nil {
		return d
	}

	d.data = nil

	var resources []Resource

	for i := range r {
		resources = append(resources, r[i])
	}

	d.dataSlice = resources

	return d
}

func (d *documentWriter) Include(r Resource) DocumentWriter {

	if d.errors != nil {
		return d
	}

	if d.Included == nil {
		d.Included = make([]Resource, 0)
	}

	d.Included = append(d.Included, r)

	return d
}

func (d *documentWriter) LinkTo(name string, value string) DocumentWriter {
	if d.Links == nil {
		d.Links = Links{}
	}

	d.Links[name] = value

	return d
}

func (d *documentWriter) SetMeta(key string, value any) DocumentWriter {
	if d.Meta == nil {
		d.Meta = Meta{}
	}

	d.Meta[key] = value

	return d
}

func (d *documentWriter) AddError(err error) DocumentWriter {
	if err == nil {
		return d
	}

	d.data = nil
	d.dataSlice = nil
	d.Included = nil

	if d.errors == nil {
		d.errors = make([]error, 0)
	}

	d.errors = append(d.errors, splitError(err)...)

	return d
}

func (d *documentWriter) Sortable() DocumentWriter {

	d.useSort = true

	return d

}

func (d *documentWriter) FieldSet() DocumentWriter {

	d.useSparseFieldSet = true

	return d

}

func (d *documentWriter) IncludeRelated() DocumentWriter {
	d.useInclude = true

	return d
}

func (d *documentWriter) WriteHeader(status int) DocumentWriter {
	d.status = status
	return d
}

func (d *documentWriter) SetHeader(key string, val string) DocumentWriter {
	if d.headers == nil {
		d.headers = map[string]string{}
	}

	d.headers[key] = val

	return d
}

func (d *documentWriter) HasErrors() bool {
	return len(d.errors) > 0
}

func (d *documentWriter) sort() error {

	if !d.useSort && d.req.URL.Query().Has("sort") {
		return Error(UnsupportedQuery, "sort")
	}

	if !d.useSort || !d.req.URL.Query().Has("sort") {
		return nil
	}

	if d.dataSlice == nil {
		return nil
	}

	sorted := make([]Resource, 0)

	for i, item := range d.dataSlice {

		x, ok := item.(Resource)
		sorted = append(sorted, x)

		if !ok {
			return fmt.Errorf("unable to parse item at data[%v]", i)
		}

	}

	sortBy := strings.Split(d.req.URL.Query()["sort"][0], ",")

	slices.SortFunc(sorted, func(a, b Resource) int {

		for _, sortByRaw := range sortBy {
			if sortByRaw == "" {
				continue
			}

			// 3. Parse the sort string (e.g., "-name")
			reverse := 1
			field := sortByRaw
			if strings.HasPrefix(sortByRaw, "-") {
				reverse = -1
				field = strings.TrimPrefix(sortByRaw, "-")
			}

			valA, _ := a.Attributes().Get(field)
			valB, _ := b.Attributes().Get(field)

			// 4. Handle comparison based on type
			var order int
			switch vA := valA.(type) {
			case string:
				vB, _ := valB.(string)
				order = strings.Compare(vA, vB)
			case float64:
				vB, _ := valB.(float64)
				order = cmp.Compare(vA, vB)
			case int:
				vB, _ := valB.(int)
				order = cmp.Compare(vA, vB)
			}

			// 5. If we found a difference, return it; otherwise, try the next field
			if order != 0 {
				return order * reverse
			}
		}

		return 0
	})

	d.dataSlice = nil
	d.AddResources(sorted)

	return nil
}

func (d *documentWriter) applySparseFieldSet() error {

	fields, err := ReadQueryParams(d.req.URL.Query())

	if err != nil {
		return err
	}

	if fields != nil {

		if d.dataSlice != nil {
			for _, r := range d.dataSlice {
				r.SetFields(fields)
			}
		}

		if d.data != nil {
			d.data.SetFields(fields)
		}

	}

	return nil
}

func (d *documentWriter) include() error {

	if !d.useInclude && d.req.URL.Query().Has("include") {
		return Error(UnsupportedQuery, "include")
	}

	if !d.useInclude || !d.req.URL.Query().Has("include") {
		return nil
	}

	data := make([]Resource, 0)
	included := make([]Resource, 0)
	query := d.req.URL.Query()["include"]
	cache := map[string]bool{}

	if d.data != nil {
		data = append(data, d.data)
	} else if d.dataSlice != nil {
		data = append(data, d.dataSlice...)
	}

	for _, _q := range query {

		rTypes := strings.Split(_q, ".")

		if len(rTypes) < 1 {
			d.AddError(Error(ForbiddenRelationship, _q))
			continue
		}

		tmpInclude := make([]Resource, 0)

		for i, _rt := range rTypes {

			tmp := data

			if i > 0 {
				tmp = tmpInclude
			}

			for _, _r := range tmp {
				res, err := GetRelatedResources(_rt, _r)

				if err != nil {
					d.AddError(err)
				} else {
					for _, r := range res {
						_, exists := cache[r.ID()] // TODO: Should probably be unique to ID + Type

						if !exists {
							tmpInclude = append(tmpInclude, r)
							cache[r.ID()] = true
						}
					}
				}
			}
		}

		included = append(included, tmpInclude...)

	}

	if d.HasErrors() {
		return nil
	}

	for _, r := range included {
		d.Include(r)
	}

	for i := range d.Included {
		r := d.Included[i]
		linkDecorator[r.Type()](d.Included[i])
	}

	return nil
}

func (d *documentWriter) MarshalJSON() ([]byte, error) {
	if d.HasErrors() {
		tmp := struct {
			Jsonapi *JsonApi `json:"jsonapi,omitempty"`
			Meta    Meta     `json:"meta,omitempty"`
			Links   Links    `json:"links,omitempty"`
			Errors  []TError `json:"errors,omitempty"`
		}{
			Jsonapi: d.Jsonapi,
			Links:   d.Links,
			Meta:    d.Meta,
		}

		tmp.Errors = make([]TError, len(d.errors))
		for i, err := range d.errors {
			if e, ok := err.(APIError); ok {
				tmp.Errors[i] = TError{
					Title:  e.Title(),
					Detail: e.Error(),
					Status: strconv.Itoa(e.Status()),
					Code:   fmt.Sprint(e.Code()),
				}
			} else {
				tmp.Errors[i] = TError{
					Detail: err.Error(),
				}
			}
		}

		return json.Marshal(tmp)

	} else if d.dataSlice != nil {
		for _, item := range d.dataSlice {
			if linkDecorator, exists := linkDecorator[item.Type()]; exists {
				linkDecorator(item)
			}
		}
		return json.Marshal(
			struct {
				Jsonapi  *JsonApi   `json:"jsonapi,omitempty"`
				Meta     Meta       `json:"meta,omitempty"`
				Links    Links      `json:"links,omitempty"`
				Data     []Resource `json:"data,omitempty"`
				Included []Resource `json:"included,omitempty"`
			}{
				Jsonapi:  d.Jsonapi,
				Links:    d.Links,
				Included: d.Included,
				Data:     d.dataSlice,
				Meta:     d.Meta,
			},
		)
	} else if d.data != nil {

		if linkDecorator, exists := linkDecorator[d.data.Type()]; exists {
			linkDecorator(d.data)
		}

		return json.Marshal(
			struct {
				Jsonapi  *JsonApi   `json:"jsonapi,omitempty"`
				Meta     Meta       `json:"meta,omitempty"`
				Links    Links      `json:"links,omitempty"`
				Data     Resource   `json:"data,omitempty"`
				Included []Resource `json:"included,omitempty"`
			}{
				Jsonapi:  d.Jsonapi,
				Links:    d.Links,
				Included: d.Included,
				Data:     d.data,
				Meta:     d.Meta,
			},
		)
	}

	return json.Marshal(
		struct {
			Jsonapi  *JsonApi   `json:"jsonapi,omitempty"`
			Meta     Meta       `json:"meta,omitempty"`
			Links    Links      `json:"links,omitempty"`
			Included []Resource `json:"included,omitempty"`
		}{
			Jsonapi:  d.Jsonapi,
			Links:    d.Links,
			Included: d.Included,
			Meta:     d.Meta,
		},
	)
}

func (d *documentWriter) Write(w http.ResponseWriter) {

	if d.writeCalled {
		log.Println("warning: write called multiple times")
		return
	}

	d.writeCalled = true

	for key, val := range d.headers {
		w.Header().Add(key, val)
	}

	w.Header().Add("Content-Type", "application/vnd.api+json")

	res := []byte("")
	var err error

	defer func() {
		res, err = json.Marshal(d)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if _, err := w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	err = d.sort()
	if err != nil {
		d.AddError(err)
	}

	if d.useSparseFieldSet {
		d.AddError(d.applySparseFieldSet())
	}

	err = d.include()
	if err != nil {
		d.AddError(err)
	}

	w.WriteHeader(d.getStatus())
}

func (d *documentWriter) getStatus() int {

	if d.status != 0 {
		return d.status
	}

	if d.HasErrors() {
		errs := make([]int, len(d.errors))

		for i := range errs {
			if e, ok := d.errors[i].(APIError); ok {
				errs[i] = e.Status()
			} else {
				errs[i] = http.StatusInternalServerError
			}
		}

		commonStatus := errs[0]
		for i := range errs {
			if errs[i] != commonStatus {
				commonStatus = 0
				break
			}
		}

		if commonStatus != 0 {
			return commonStatus
		}

		gStatus := []int{0, 0}
		for i := range errs {
			if errs[i] >= 400 && errs[i] < 500 {
				gStatus[0]++
			}
			if errs[i] >= 500 && errs[i] < 600 {
				gStatus[1]++
			}
		}

		if gStatus[0] > gStatus[1] {
			return 400
		}
		if gStatus[1] > gStatus[0] {
			return 500
		}

		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func splitError(err error) []error {
	errout := make([]error, 0)

	e, ok := err.(ErrorGroup)

	if !ok {
		return append(errout, err)
	}

	for _, nerr := range e.Split() {
		errout = append(errout, splitError(nerr)...)
	}

	return errout
}

func GetRelatedResources(relType string, r Resource) ([]Resource, error) {

	items := make([]Resource, 0)

	rels := r.Relationships().Get(relType)
	if rels == nil {
		return items, nil
	}

	for _, rel := range rels.Data() {

		getter, exists := Getters[rel.Type]

		if !exists {
			return nil, (Error(MissingGetter, rel.Type))
		}

		resource, err := getter(rel.ID)

		if err != nil {
			return nil, err
		} else {
			items = append(items, resource)
		}

	}

	return items, nil

}
