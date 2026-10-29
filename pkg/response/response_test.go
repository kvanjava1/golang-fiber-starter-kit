package response_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-starter-kit/pkg/response"
)

func TestResponse_Success(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	res := response.Success("Success retrieved", data)

	assert.True(t, res.Success)
	assert.Equal(t, "Success retrieved", res.Message)
	assert.Equal(t, data, res.Data)
	assert.NotEmpty(t, res.Meta.Timestamp)
}

func TestResponse_Error(t *testing.T) {
	res := response.Error("Bad request", nil)

	assert.False(t, res.Success)
	assert.Equal(t, "Bad request", res.Message)
	assert.Nil(t, res.Data)
	assert.NotEmpty(t, res.Meta.Timestamp)
}

func TestResponse_ValidationFailed(t *testing.T) {
	fieldErrors := []response.FieldError{
		{Field: "email", Message: "email is required"},
	}
	res := response.ValidationFailed(fieldErrors)

	assert.False(t, res.Success)
	assert.Equal(t, "Validation failed", res.Message)
	assert.Len(t, res.Errors, 1)
	assert.Equal(t, "email", res.Errors[0].Field)
}

func TestResponse_MapError(t *testing.T) {
	tests := []struct {
		err          error
		expectedCode int
		expectedMsg  string
	}{
		{response.ErrNotFound, http.StatusNotFound, "Resource not found"},
		{response.ErrConflict, http.StatusConflict, "Resource already exists"},
		{response.ErrUnauthorized, http.StatusUnauthorized, "Unauthorized access"},
		{response.ErrForbidden, http.StatusForbidden, "Access forbidden"},
		{response.ErrValidation, http.StatusBadRequest, "Validation failed"},
		{errors.New("random db error"), http.StatusInternalServerError, "Internal server error"},
	}

	for _, tt := range tests {
		code, msg := response.MapError(tt.err)
		assert.Equal(t, tt.expectedCode, code)
		assert.Equal(t, tt.expectedMsg, msg)
	}
}
