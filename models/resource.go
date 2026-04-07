package models

import (
	"encoding/json"
)

type Resource interface {
	ID() string
	Type() string
	Attributes() Attributes
	Relationships() Relationships
	Links() *Links
	Meta() *Meta

	SetFields(map[string][]string)
}

func NewResource(ID string, Type string) Resource {
	return &resource{
		id:            ID,
		_type:         Type,
		attributes:    new(attributes),
		relationships: new(relationships),
		links:         Links{},
		meta:          Meta{},
	}
}

type resource struct {
	id            string
	_type         string
	attributes    Attributes
	relationships Relationships
	links         Links
	meta          Meta

	fields []string
}

func (r *resource) ID() string                   { return r.id }
func (r *resource) Type() string                 { return r._type }
func (r *resource) Attributes() Attributes       { return r.attributes }
func (r *resource) Relationships() Relationships { return r.relationships }
func (r *resource) Links() *Links                { return &r.links }
func (r *resource) Meta() *Meta                  { return &r.meta }

func (d *resource) MarshalJSON() ([]byte, error) {
	r := struct {
		ID            string      `json:"id"`
		Type          string      `json:"type"`
		Attributes    interface{} `json:"attributes,omitempty"`
		Relationships interface{} `json:"relationships,omitempty"`
		Links         interface{} `json:"links,omitempty"`
		Meta          interface{} `json:"meta,omitempty"`
	}{
		ID:   d.id,
		Type: d._type,
	}

	if d.attributes != nil && !d.attributes.IsNil() {

		if d.fields != nil {
			filteredAttrs := NewAttributes()
			hasData := false
			for _, f := range d.fields {
				if val, exists := d.attributes.Get(f); exists {
					filteredAttrs.Set(f, val)
					hasData = true
				}
			}
			if hasData {
				r.Attributes = filteredAttrs
			}
		} else {
			r.Attributes = d.attributes
		}

	}

	if d.relationships != nil && !d.relationships.IsNil() {
		r.Relationships = d.relationships
	}

	if len(d.meta) > 0 {
		r.Meta = d.meta
	}

	if len(d.links) > 0 {
		r.Links = d.links
	}

	return json.Marshal(r)

}

func (r *resource) UnmarshalJSON(src []byte) error {
	obj := struct {
		ID            string        `json:"id"`
		Type          string        `json:"type"`
		Attributes    attributes    `json:"attributes,omitempty"`
		Relationships relationships `json:"relationships,omitempty"`
		Links         Links         `json:"links,omitempty"`
		Meta          Meta          `json:"meta,omitempty"`
	}{}

	err := json.Unmarshal(src, &obj)

	if err != nil {
		return err
	}

	r.id = obj.ID
	r._type = obj.Type
	r.attributes = &obj.Attributes
	r.relationships = &obj.Relationships
	r.links = obj.Links
	r.meta = obj.Meta

	return nil
}

func (r *resource) SetFields(filter map[string][]string) {

	for dataType, fields := range filter {

		if dataType == r._type {

			if r.fields == nil {
				r.fields = make([]string, 0)
			}

			r.fields = append(r.fields, fields...)
		}

	}

}

func NewBasicResource(ID string, Type string) Resource {
	return &basicResource{id: ID, _type: Type}
}

type basicResource struct {
	id    string
	_type string
}

func (r *basicResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}{ID: r.id, Type: r._type})
}

func (r *basicResource) ID() string                    { return r.id }
func (r *basicResource) Type() string                  { return r._type }
func (r *basicResource) Attributes() Attributes        { return NewAttributes() }
func (r *basicResource) Relationships() Relationships  { return make(relationshipsStub) }
func (r *basicResource) Links() *Links                 { return &Links{} }
func (r *basicResource) Meta() *Meta                   { return &Meta{} }
func (r *basicResource) SetFields(map[string][]string) {}

type relationshipsStub map[string]Relationship

func (relationshipsStub) Add(string, Relationship)        {}
func (relationshipsStub) Get(string) Relationship         { return NewRelationship() }
func (relationshipsStub) Each(func(string, Relationship)) {}
func (relationshipsStub) Len() int                        { return 0 }
func (relationshipsStub) IsNil() bool                     { return true }
