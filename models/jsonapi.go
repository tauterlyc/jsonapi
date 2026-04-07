package models

type JsonApi struct {
	Version  string                 `json:"version,omitempty"`
	Ext      []string               `json:"ext,omitempty"`
	Profiles []string               `json:"profiles,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}
