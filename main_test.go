package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/ChrisTheShark/simple-admission-controller"
	"github.com/stretchr/testify/assert"
)

func TestValidateNoBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/validate", nil)
	w := httptest.NewRecorder()

	Validate(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestValidateIncorrectBody(t *testing.T) {
	user := struct {
		FirstName string
	}{
		"Chuck",
	}

	bs, _ := json.Marshal(&user)
	r := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(bs))
	w := httptest.NewRecorder()

	Validate(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
