package config

import (
	"errors"

	"github.com/gitamix/lint/config"
)

// Validation checks a gitamix lint configuration for completeness.
type Validation struct {
	// cfg stores the lint configuration that gets validated.
	cfg config.Config
}

// NewValidation creates a new Validation instance
// with the provided lint configuration.
func NewValidation(cfg config.Config) Validation {
	return Validation{
		cfg: cfg,
	}
}

// Validate reports whether the lint configuration is complete.
//
// It returns a non-nil error when the default branch name
// in the configuration is empty; otherwise it returns nil.
func (v Validation) Validate() error {
	if v.cfg.Branch().Default().Name() == "" {
		return errors.New("empty default branch name")
	}
	return nil
}
