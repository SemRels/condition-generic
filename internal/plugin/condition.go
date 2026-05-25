package plugin

import (
	"fmt"
	"os"
	"strings"
)

type Condition struct {
	env func(string) string
}

func New() *Condition { return &Condition{env: os.Getenv} }

func NewWithEnv(env func(string) string) *Condition { return &Condition{env: env} }

func (c *Condition) Check() error {
	name := strings.TrimSpace(c.env("SEMREL_PLUGIN_ENV_VAR"))
	if name == "" {
		return nil
	}

	got := c.env(name)
	want := c.env("SEMREL_PLUGIN_ENV_VALUE")
	if got != want {
		return fmt.Errorf("environment mismatch: %s=%q, want %q", name, got, want)
	}

	return nil
}
