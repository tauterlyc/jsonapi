package models

type Attributes interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	IsNil() bool
}

func NewAttributes() Attributes {
	return &attributes{}
}

type attributes map[string]interface{}

func (m *attributes) Set(key string, value interface{}) {
	if *m == nil {
		*m = make(attributes)
	}

	(*m)[key] = value
}

func (m *attributes) Get(key string) (interface{}, bool) {
	if *m == nil {
		return nil, false
	}

	value, exists := (*m)[key]

	return value, exists
}

func (m *attributes) IsNil() bool {
	return *m == nil
}
