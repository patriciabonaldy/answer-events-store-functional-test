package domain_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/patriciabonaldy/answer-events-store-functional/domain"
)

func TestRouter_ErrorResponse(t *testing.T) {
	err := domain.ErrorResponse{
		ErrorCode:    "01",
		ErrorMessage: "fail",
	}

	assert.Equal(t, "fail", err.Error())
}

func TestRouter_ErrorHttp(t *testing.T) {
	err400 := domain.ErrorHttp{
		Cause:          errors.New("chain error"),
		Message:        "bad request",
		ExternalStatus: "01",
		HTTPStatus:     400,
	}

	err500 := domain.ErrorHttp{
		Cause:          nil,
		Message:        "internal server error",
		ExternalStatus: "02",
		HTTPStatus:     500,
	}

	assert.Equal(t, "bad request: chain error", err400.Error())
	assert.Equal(t, "internal server error", err500.Error())
}
