package models

type Meta map[string]interface{}

func (m *Meta) Add(key string, value interface{}) {
	if *m == nil {
		*m = make(Meta)
	}

	(*m)[key] = value
}

func (m *Meta) Get(key string) interface{} {
	if *m == nil {
		return nil
	}

	return (*m)[key]
}
