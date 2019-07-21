package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/ChrisTheShark/simple-admission-controller"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
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

func TestValidateValid(t *testing.T) {
	podJson := `
	{
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": {
		   "name": "compliant-pod",
		   "namespace": "compliant"
		},
		"spec": {
		   "containers": [
			  {
				 "image": "nginx:1.7.9",
				 "name": "compliant-pod",
				 "ports": [
					{
					   "containerPort": 80
					}
				 ]
			  }
		   ]
		}
	 }
	`
	arRequest := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: []byte(podJson),
			},
		},
	}

	bs, _ := json.Marshal(&arRequest)
	r := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(bs))
	w := httptest.NewRecorder()

	Validate(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestValidateInvalid(t *testing.T) {
	podJson := `
	{
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": {
		   "name": "noncompliant-pod",
		   "namespace": "default"
		},
		"spec": {
		   "containers": [
			  {
				 "image": "nginx:1.7.9",
				 "name": "noncompliant-pod",
				 "ports": [
					{
					   "containerPort": 80
					}
				 ]
			  }
		   ]
		}
	 }
	`
	arRequest := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: []byte(podJson),
			},
		},
	}

	bs, _ := json.Marshal(&arRequest)
	r := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(bs))
	w := httptest.NewRecorder()

	Validate(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	bs, _ = ioutil.ReadAll(resp.Body)
	respObject := v1beta1.AdmissionReview{}
	json.Unmarshal(bs, &respObject)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, false, respObject.Response.Allowed)
}
