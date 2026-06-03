// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package plugin

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Condition struct {
	env func(string) string
}

func New() *Condition { return &Condition{env: os.Getenv} }

func NewWithEnv(env func(string) string) *Condition { return &Condition{env: env} }

func (c *Condition) Check() error {
	if commands := c.env("SEMREL_PLUGIN_COMMAND"); commands != "" {
		for _, command := range strings.Split(commands, "\n") {
			command = strings.TrimSpace(command)
			if command == "" {
				continue
			}

			cmd := exec.Command("sh", "-c", command)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				msg := strings.TrimSpace(stderr.String())
				if msg != "" {
					return fmt.Errorf("command failed: %s: %w: %s", command, err, msg)
				}
				return fmt.Errorf("command failed: %s: %w", command, err)
			}
		}

		return nil
	}

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
