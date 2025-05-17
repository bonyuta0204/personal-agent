package main

import "github.com/bonyuta0204/personal-agent/go/config"

// AppContext holds the application context including configuration
type AppContext struct {
	Config *config.Config
}

// NewAppContext creates a new application context with the given configuration
func NewAppContext(cfg *config.Config) *AppContext {
	return &AppContext{
		Config: cfg,
	}
}
