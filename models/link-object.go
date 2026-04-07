package models

type LinkObject struct {
	Href        string `json:"href"`
	Rel         string `json:"rel,omitempty"`
	DescribedBy string `json:"describedBy,omitempty"`
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
}

type Links map[string]string

func (m *Links) Add(key string, href string) {
	if *m == nil {
		*m = make(Links)
	}

	(*m)[key] = href
}

// func (m *LinksMap) UnmarshalJSON(src []byte) error {
// 	var rawMap map[string]json.RawMessage
// 	if err := json.Unmarshal(src, &rawMap); err != nil {
// 		return err
// 	}

// 	*m = make(LinksMap)

// 	for key, rawValue := range rawMap {
// 		var obj LinkObject

// 		if err := json.Unmarshal(rawValue, obj); err != nil {
// 			return err
// 		}
// 		(*m)[key] = obj
// 	}
// 	return nil
// }
