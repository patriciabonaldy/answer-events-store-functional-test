package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/patriciabonaldy/answer-events-store-functional/domain"
)

func TestMapGetValidString(t *testing.T) {
	input := map[string]interface{}{
		"l1": map[string]interface{}{
			"v1": "teste1",
		},
	}

	fields := domain.NewFieldsFromMap(input)
	assert.Equal(t, "teste1", fields.GetString("l1.v1"))
}

func TestMapGetInvalidString(t *testing.T) {
	input := map[string]interface{}{
		"l1": map[string]interface{}{
			"v1": nil,
		},
	}

	fields := domain.NewFieldsFromMap(input)
	assert.Empty(t, fields.GetString("l1.v1"))
}

func TestMapGetFields(t *testing.T) {
	input := map[string]interface{}{
		"l1": map[string]interface{}{
			"v1": "teste1",
			"v2": "teste2",
		},
	}

	fields := domain.NewFieldsFromMap(input)
	recv := fields.GetFields("l1")

	assert.Equal(t, "teste1", recv.GetString("v1"))
	assert.Equal(t, "teste2", recv.GetString("v2"))
}

func TestMapGetInt_WhenJsonNumber(t *testing.T) {
	var number json.Number = "10"

	input := map[string]interface{}{
		"l1": map[string]interface{}{
			"v1": number,
		},
	}

	fields := domain.NewFieldsFromMap(input)
	assert.Equal(t, int64(10), fields.GetInt("l1.v1"))
}

func TestMapGetInt_WhenInteger(t *testing.T) {
	input := map[string]interface{}{
		"l1": map[string]interface{}{
			"v1": 10,
		},
	}

	fields := domain.NewFieldsFromMap(input)
	assert.Equal(t, int64(10), fields.GetInt("l1.v1"))
}

func TestMapGetInt_WhenFloat(t *testing.T) {
	input := map[string]interface{}{
		"l1": map[string]interface{}{
			"v1": 10.5,
		},
	}

	fields := domain.NewFieldsFromMap(input)
	assert.Equal(t, int64(10), fields.GetInt("l1.v1"))
}

func TestMapGetInvalidInt(t *testing.T) {
	input := map[string]interface{}{
		"l1": map[string]interface{}{
			"v1": nil,
		},
	}

	fields := domain.NewFieldsFromMap(input)
	assert.Equal(t, int64(0), fields.GetInt("l1.v1"))
}

func TestMapSetString_WhenMapExists(t *testing.T) {
	input := map[string]interface{}{
		"l1": map[string]interface{}{
			"v1": "teste1",
		},
	}

	fields := domain.NewFieldsFromMap(input)
	fields.Set("l1.v1", "teste2")

	assert.Equal(t, "teste2", fields.GetString("l1.v1"))
}

func TestMapSetString_WhenMapNotExists(t *testing.T) {
	input := map[string]interface{}{}

	fields := domain.NewFieldsFromMap(input)
	fields.Set("l1.v1", "teste2")

	assert.Equal(t, "teste2", fields.GetString("l1.v1"))
}

func TestInvalidGet(t *testing.T) {
	input := map[string]interface{}{}
	fields := domain.NewFieldsFromMap(input)

	assert.Nil(t, fields.Get("l1.v1"))
}
