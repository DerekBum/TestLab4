package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	err := RespondJSON(w, 200, nil)
	assert.Nil(t, err)
}

func TestServerError(t *testing.T) {
	w := httptest.NewRecorder()
	ServerError(w)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequest(w, "lol")
	assert.Equal(t, w.Code, http.StatusBadRequest)
}
