package httpresp

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	var (
		body = map[string]any{
			"foo": "bar",
		}
		statusCode = http.StatusOK
	)
	response := New(statusCode, body)
	assert.Equal(t, body, response.body)
	assert.Equal(t, statusCode, response.statusCode)
}

func TestNewOK(t *testing.T) {
	body := map[string]any{
		"foo": "bar",
	}
	response := NewOK(body)
	assert.Equal(t, body, response.body)
	assert.Equal(t, http.StatusOK, response.statusCode)
}

func TestRenderJSON(t *testing.T) {
	r := httptest.NewRecorder()

	payload := map[string]any{
		"foo": "bar",
	}
	renderJSON(r, http.StatusOK, &payload)
	assert.Equal(t, r.Code, http.StatusOK)
	assert.Equal(t, r.Header().Get("Content-Type"), "application/json")

	responsePayload := make(map[string]any)
	err := json.NewDecoder(r.Body).Decode(&responsePayload)
	assert.NoError(t, err)

	if assert.NotEmpty(t, responsePayload) {
		assert.Equal(t, payload, responsePayload)
	}
}

func TestRender(t *testing.T) {
	r := httptest.NewRecorder()

	payload := map[string]any{
		"foo": "bar",
	}
	response := New(http.StatusOK, &payload)
	Render(r, response)
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "application/json", r.Header().Get("Content-Type"))

	responsePayload := make(map[string]any)
	err := json.NewDecoder(r.Body).Decode(&responsePayload)
	assert.NoError(t, err)

	if assert.NotEmpty(t, responsePayload) {
		assert.Equal(t, payload, responsePayload)
	}
}
