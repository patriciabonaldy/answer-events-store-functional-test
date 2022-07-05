package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Fields stores all fields from json as an interface
type Fields map[string]interface{}

// NewFieldsFromMap initialize fields interface
func NewFieldsFromMap(fields map[string]interface{}) Fields {
	for k, v := range fields {
		if m, ok := v.(map[string]interface{}); ok {
			fields[k] = NewFieldsFromMap(m)
		}
	}

	return fields
}

// GetString returns the value from path as string
func (ref Fields) GetString(path string) string {
	v := ref.Get(path)
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}

// GetFields returns all values from path
func (ref Fields) GetFields(path string) Fields {
	raw := ref.Get(path)
	m, ok := raw.(map[string]interface{})

	if !ok {
		m, _ = raw.(Fields)
	}

	return m
}

// GetInt returns the value from path as int64
func (ref Fields) GetInt(path string) int64 {
	raw := ref.Get(path)
	switch m := raw.(type) {
	case json.Number:
		v, _ := m.Int64()
		return v
	case int:
		return int64(m)
	case float64:
		return int64(m)
	}

	return 0
}

// Get returns all child fields from path
func (ref Fields) Get(path string) interface{} {
	fields := strings.Split(path, ".")
	if len(fields) == 0 {
		return nil
	}

	field := fields[0]
	if len(fields) == 1 {
		return ref[field]
	}

	m, ok := ref[field].(map[string]interface{})
	if !ok {
		m, _ = ref[field].(Fields)
	}

	if m == nil {
		return nil
	}

	innerPath := strings.Join(fields[1:], ".")
	return NewFieldsFromMap(m).Get(innerPath)
}

// Set changes the value of a path variable
func (ref Fields) Set(path string, value interface{}) {
	fields := strings.Split(path, ".")
	if len(fields) == 0 {
		return
	}

	field := fields[0]
	if len(fields) == 1 {
		ref[field] = value
		return
	}

	m, ok := ref[field].(map[string]interface{})
	if !ok {
		m, _ = ref[field].(Fields)
	}

	if m == nil {
		m = Fields{}
		ref[field] = m
	}

	innerPath := strings.Join(fields[1:], ".")
	NewFieldsFromMap(m).Set(innerPath, value)
}
