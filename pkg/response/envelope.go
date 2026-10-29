package response

import "time"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Meta struct {
	RequestID string `json:"requestId,omitempty"`
	Timestamp string `json:"timestamp"`
	Page      int    `json:"page,omitempty"`
	PerPage   int    `json:"perPage,omitempty"`
	Total     int64  `json:"total,omitempty"`
}

func NewMeta(requestID string) Meta {
	return Meta{
		RequestID: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

type Envelope struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    any          `json:"data,omitempty"`
	Errors  []FieldError `json:"errors,omitempty"`
	Meta    Meta         `json:"meta"`
}
