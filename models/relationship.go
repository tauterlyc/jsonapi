package models

import (
	"encoding/json"
	"errors"
)

type Relationship interface {
	Links() *Links
	Meta() *Meta
	Data() []ResourceIdentifier

	SetData([]ResourceIdentifier)

	LinkTo(name string, url string)
	GetResourceIdentifiers() ([]ResourceIdentifier, bool)
}

type relationship struct {
	links Links
	meta  Meta

	data      *ResourceIdentifier
	dataSlice []ResourceIdentifier
}

func NewRelationship() Relationship {
	return &relationship{}
}

func (r *relationship) Links() *Links {
	return &r.links
}

func (r *relationship) Meta() *Meta {
	return &r.meta
}

func (r *relationship) Data() []ResourceIdentifier {
	data, _ := r.GetResourceIdentifiers()
	return data
}

func (r *relationship) SetData(re []ResourceIdentifier) {
	r.data = nil
	if r.dataSlice == nil {
		r.dataSlice = make([]ResourceIdentifier, 0)
	}

	r.dataSlice = append(r.dataSlice, re...)
}

func (r *relationship) LinkTo(name string, url string) {
	if r.links == nil {
		r.links = make(Links)
	}

	r.links[name] = url
}

func (r *relationship) GetResourceIdentifiers() ([]ResourceIdentifier, bool) {
	if r.data != nil {
		return []ResourceIdentifier{*r.data}, true
	}

	if r.dataSlice != nil {
		return r.dataSlice, false
	}

	return nil, false
}

func (r *relationship) MarshalJSON() ([]byte, error) {
	if r.data != nil {

		return json.Marshal(struct {
			Links Links                  `json:"links,omitempty"`
			Meta  map[string]interface{} `json:"meta,omitempty"`
			Data  *ResourceIdentifier    `json:"data,omitempty"`
		}{
			Links: r.links,
			Meta:  r.meta,
			Data:  r.data,
		})

	} else if r.dataSlice != nil {

		return json.Marshal(struct {
			Links Links                  `json:"links,omitempty"`
			Meta  map[string]interface{} `json:"meta,omitempty"`
			Data  []ResourceIdentifier   `json:"data,omitempty"`
		}{
			Links: r.links,
			Meta:  r.meta,
			Data:  r.dataSlice,
		})

	}

	return json.Marshal(struct {
		Links Links                  `json:"links,omitempty"`
		Meta  map[string]interface{} `json:"meta,omitempty"`
	}{
		Links: r.links,
		Meta:  r.meta,
	})
}

func (r *relationship) UnmarshalJSON(src []byte) error {
	obj := struct {
		Links Links                  `json:"links,omitempty"`
		Meta  map[string]interface{} `json:"meta,omitempty"`
	}{}

	err := json.Unmarshal(src, &obj)

	if err != nil {
		return err
	}

	r.links = obj.Links
	r.meta = obj.Meta

	objS := struct {
		Data *ResourceIdentifier `json:"data,omitempty"`
	}{}
	sErr := json.Unmarshal(src, &objS)

	if sErr == nil {
		r.data = objS.Data
		return nil
	}

	objM := struct {
		Data []ResourceIdentifier `json:"data,omitempty"`
	}{}
	mErr := json.Unmarshal(src, &objM)

	if mErr == nil {
		r.dataSlice = objM.Data
		return nil
	}

	return errors.Join(err, sErr, mErr)
}