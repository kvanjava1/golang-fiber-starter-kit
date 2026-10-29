package response

import (
	"time"
)

// Success membangun format response JSON standar sukses
func Success(message string, data any, meta ...Meta) Envelope {
	m := Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)}
	if len(meta) > 0 {
		m = meta[0]
		if m.Timestamp == "" {
			m.Timestamp = time.Now().UTC().Format(time.RFC3339)
		}
	}

	return Envelope{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    m,
	}
}

// Error membangun format response JSON error
func Error(message string, errs []FieldError, meta ...Meta) Envelope {
	m := Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)}
	if len(meta) > 0 {
		m = meta[0]
		if m.Timestamp == "" {
			m.Timestamp = time.Now().UTC().Format(time.RFC3339)
		}
	}

	return Envelope{
		Success: false,
		Message: message,
		Errors:  errs,
		Meta:    m,
	}
}

// ValidationFailed membangun format response JSON khusus error validasi form / DTO
func ValidationFailed(errors []FieldError, meta ...Meta) Envelope {
	return Error("Validation failed", errors, meta...)
}
