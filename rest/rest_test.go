package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/patriciabonaldy/answer-events-store-functional/domain"
	"github.com/patriciabonaldy/answer-events-store-functional/rest"
)

func successHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"acquirer_transaction_id": "abbbbbbbbbccccddddddddfffff"}`))
	if err != nil {
		panic("Problem")
	}
}

func errorResponseHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(`{"error_code": "01", "error_message": "fail"}`))
	if err != nil {
		panic("Problem")
	}
}

func errorHttpHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(`{"error": "01", "message": "bad request", "status": 400}`))
	if err != nil {
		panic("Problem")
	}
}

func errorGenericHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(errors.New("generic failure").Error()))
	if err != nil {
		panic("Problem")
	}
}

func TestHttpRequest_Success(t *testing.T) {
	response := domain.Fields{}

	server := httptest.NewServer(http.HandlerFunc(successHandleFunc))
	err := rest.HttpRequest(server.URL, "POST", nil, &response)

	assert.NoError(t, err)
	assert.NotEmpty(t, response.GetString("acquirer_transaction_id"))
}

func TestHttpRequest_ErrorResponse(t *testing.T) {
	var errResp error

	server := httptest.NewServer(http.HandlerFunc(errorResponseHandleFunc))
	err := rest.HttpRequest(server.URL, "POST", nil, &errResp)

	assert.Error(t, err)
}

func TestHttpRequest_ErrorHttp(t *testing.T) {
	var errHttp error

	server := httptest.NewServer(http.HandlerFunc(errorHttpHandleFunc))
	err := rest.HttpRequest(server.URL, "POST", nil, &errHttp)

	assert.Error(t, err)
}

func TestHttpRequest_ErrorGeneric(t *testing.T) {
	var errGeneric error

	server := httptest.NewServer(http.HandlerFunc(errorGenericHandleFunc))
	err := rest.HttpRequest(server.URL, "POST", nil, &errGeneric)

	assert.Error(t, err)
}

func TestRest_Success(t *testing.T) {
	request := domain.Fields{}
	request["param"] = "value"

	server := httptest.NewServer(http.HandlerFunc(successHandleFunc))
	response, err := rest.Rest(server.URL, "POST", request)

	assert.NoError(t, err)
	assert.NotEmpty(t, response.GetString("acquirer_transaction_id"))
}
