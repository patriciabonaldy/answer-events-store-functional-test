package config

import (
	"fmt"
	"os"
)

// Environment Struct
type Environment struct {
	RouterURL string
}

// Environments Type
type Environments map[string]Environment

var (
	environments = map[string]Environment{
		"": { // Run tests in local machine
			RouterURL: "http://localhost:8080",
		},
		"staging": { // Run tests in staging
			RouterURL: "https://domain_url/v1",
		},
	}
)

// MustGetEnvironment returns the environment scope
func MustGetEnvironment() Environment {
	scope := os.Getenv("SCOPE")
	env, exists := environments[scope]
	if !exists {
		panic(fmt.Errorf("environment [%s] not found", scope))
	}

	return env
}
