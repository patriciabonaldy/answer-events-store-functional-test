package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"

	"github.com/patriciabonaldy/answer-events-store-functional/domain"
)

func HttpRequest(url string, method string, request interface{}, response interface{}) error {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	fmt.Fprintf(os.Stderr, "_request: url(%s): %s\n", url, reqBody)

	reqHTTP, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	if reqBody != nil {
		reqHTTP.Header.Set("Content-Type", "application/json")
	}
	reqHTTP.Header.Set("x-request-id", newRequestID())
	reqHTTP.Header.Set("x-flow-starter", "true")

	var client = &http.Client{
		Timeout: time.Second * 100,
	}

	resHTTP, err := client.Do(reqHTTP)
	if err != nil {
		return fmt.Errorf("http.Do: %w", err)
	}
	defer resHTTP.Body.Close()

	resBody, err := ioutil.ReadAll(resHTTP.Body)
	if err != nil && err != io.EOF {
		return fmt.Errorf("io.Read: %w, Payload: %s", err, resBody)
	}

	fmt.Fprintf(os.Stderr, "_response: url(%s): headers: %+v: %s\n", url, resHTTP.Header, resBody)
	return setResponse(resBody, resHTTP, response)
}

func setResponse(resBody []byte, resHTTP *http.Response, response interface{}) error {
	// set transaction valid response
	if resHTTP.StatusCode >= 200 && resHTTP.StatusCode <= 299 {
		err := json.Unmarshal(resBody, &response)
		if err == nil {
			return nil
		}
	}

	// set transaction error response
	errResp := domain.ErrorResponse{}
	err := json.Unmarshal(resBody, &errResp)
	if err == nil && errResp.ErrorMessage != "" {
		return domain.ErrorHttp{
			Cause:          errResp,
			Message:        string(resBody),
			ExternalStatus: resHTTP.Status,
			HTTPStatus:     resHTTP.StatusCode,
		}
	}

	// set http error response
	errHttp := domain.ErrorHttp{}
	err = json.Unmarshal(resBody, &errHttp)
	if err == nil && errHttp.Message != "" {
		return errHttp
	}

	// set custom error response
	return domain.ErrorHttp{
		Message:        string(resBody),
		ExternalStatus: resHTTP.Status,
		HTTPStatus:     resHTTP.StatusCode,
	}
}

func Rest(url string, method string, reqFields domain.Fields) (resFields domain.Fields, err error) {
	err = restAndDecodeTo(url, method, reqFields, &resFields)
	return domain.NewFieldsFromMap(resFields), err
}

func restAndDecodeTo(url string, method string, reqFields domain.Fields, decoded interface{}) (err error) {
	var reqBody []byte
	if len(reqFields) > 0 {
		reqBody, err = json.Marshal(reqFields)
		if err != nil {
			return fmt.Errorf("json.Marshal: %w", err)
		}

		fmt.Fprintf(os.Stderr, "_request: url(%s): %s\n", url, reqBody)
	}

	reqHTTP, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	if reqBody != nil {
		reqHTTP.Header.Set("Content-Type", "application/json")
	}

	var client = &http.Client{
		Timeout: time.Second * 100,
	}

	resHTTP, err := client.Do(reqHTTP)
	if err != nil {
		return fmt.Errorf("http.Do: %w", err)
	}
	defer resHTTP.Body.Close()

	resBody, err := ioutil.ReadAll(resHTTP.Body)
	if err != nil && err != io.EOF {
		return fmt.Errorf("io.Read: %w, Payload: %s", err, resBody)
	}

	fmt.Fprintf(os.Stderr, "_response: url(%s): headers: %+v: %s\n", url, resHTTP.Header, resBody)

	if resHTTP.StatusCode < 200 || resHTTP.StatusCode > 499 {
		return fmt.Errorf("http error: %d - %s, Payload %s", resHTTP.StatusCode, resHTTP.Status, resBody)
	}

	if len(resBody) > 0 {
		return DecodeMap(resBody, decoded)
	}

	return nil
}

func DecodeMap(resBody []byte, decoded interface{}) error {
	dec := json.NewDecoder(bytes.NewBuffer(resBody))
	dec.UseNumber()

	err := dec.Decode(decoded)
	if err != nil {
		return fmt.Errorf("json.Decode: %+v", err)
	}
	return nil
}

func newRequestID() string {
	u, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return u.String()
}
