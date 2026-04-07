package models

type TError struct {
	ID     string                 `json:"id,omitempty"`
	Links  *ErrorLinks            `json:"links,omitempty"`
	Status string                 `json:"status,omitempty"`
	Code   string                 `json:"code,omitempty"`
	Title  string                 `json:"title,omitempty"`
	Detail string                 `json:"detail,omitempty"`
	Source *ErrorSource           `json:"source,omitempty"`
	Meta   map[string]interface{} `json:"Meta,omitempty"`
}

type ErrorLinks struct {
	// A link that leads to further details about the particular occurrence of the problem.
	About string `json:"about,omitempty"`
	// A link identifying the type of error, which should return a general explanation of the error.
	Type string `json:"type,omitempty"`
}

// Source object that contains references to the primary source of the error.
type ErrorSource struct {
	// A JSON Pointer to the value in the request document that caused the error.
	Pointer string `json:"pointer,omitempty"`
	// A string indicating which URI query parameter caused the error.
	Parameter string `json:"parameter,omitempty"`
	// A string indicating the name of a single request header that caused the error.
	Header string `json:"header,omitempty"`
}
