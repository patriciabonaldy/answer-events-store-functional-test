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
		"": { // Send tests from local machine to staging
			RouterURL: "https://domain_url/v1",
		},
		"local": { // Send tests from local machine to local machine services
			RouterURL: "http://localhost:8082",
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
