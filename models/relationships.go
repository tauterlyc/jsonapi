package models

import "encoding/json"

type Relationships interface {
	Add(key string, rel Relationship)
	Get(key string) Relationship
	Each(func(name string, relationship Relationship))
	Len() int
	IsNil() bool
}

func NewRelationships() Relationships {
	return &relationships{}
}

type relationships map[string]Relationship

func (m *relationships) Add(key string, rel Relationship) {
	if *m == nil {
		(*m) = make(relationships)
	}

	if (*m)[key] == nil {
		(*m)[key] = rel
	}
}

func (m *relationships) Get(key string) Relationship {
	if *m == nil {
		return nil
	}

	return (*m)[key]
}

func (m *relationships) Each(fn func(string, Relationship)) {
	for name, v := range *m {
		fn(name, v)
	}
}

func (m *relationships) Len() int {
	if *m == nil {
		return 0
	}

	return len(*m)
}

func (m *relationships) IsNil() bool {
	return *m == nil
}

func (m *relationships) UnmarshalJSON(src []byte) error {
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(src, &rawMap); err != nil {
		return err
	}

	*m = make(relationships)

	for key, rawValue := range rawMap {
		obj := NewRelationship()

		if err := json.Unmarshal(rawValue, obj); err != nil {
			return err
		}
		(*m)[key] = obj
	}
	return nil
}
