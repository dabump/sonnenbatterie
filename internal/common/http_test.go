package common

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalServerError(t *testing.T) {
	reqRecorder := httptest.NewRecorder()
	InternalServerError(reqRecorder, errors.New("my-error"))

	assert.Equal(t, http.StatusInternalServerError, reqRecorder.Result().StatusCode)
	assert.Equal(t, "application/text", reqRecorder.Header().Get("Content-Type"))

	bodyBytes, err := io.ReadAll(reqRecorder.Body)
	assert.NoError(t, err)
	assert.Equal(t, "my-error", string(bodyBytes))
}

func TestTooManyRequests(t *testing.T) {
	reqRecorder := httptest.NewRecorder()
	TooManyRequests(reqRecorder)

	assert.Equal(t, http.StatusTooManyRequests, reqRecorder.Result().StatusCode)
	assert.Equal(t, "application/text", reqRecorder.Header().Get("Content-Type"))
}
