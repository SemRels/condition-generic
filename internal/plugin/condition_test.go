// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package plugin

import (
	"strings"
	"testing"
)

func env(kv map[string]string) func(string) string {
	return func(key string) string { return kv[key] }
}

func TestCheck_PermissiveWithoutConfig(t *testing.T) {
	t.Parallel()

	if err := NewWithEnv(env(map[string]string{})).Check(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_CommandSuccess(t *testing.T) {
	t.Parallel()

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND": "true",
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_CommandFailure(t *testing.T) {
	t.Parallel()

	command := "printf 'boom\\n' >&2; exit 1"
	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND": command,
	})).Check()
	if err == nil {
		t.Fatal("expected command failure")
	}
	if !strings.Contains(err.Error(), command) {
		t.Fatalf("expected failing command in error, got %v", err)
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected stderr in error, got %v", err)
	}
}

func TestCheck_MultipleCommandsSuccess(t *testing.T) {
	t.Parallel()

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND": "true\n:",
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_MultipleCommandsFailureStopsAtFailingCommand(t *testing.T) {
	t.Parallel()

	failing := "printf 'stop\\n' >&2; exit 2"
	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND": "true\n" + failing + "\ntrue",
	})).Check()
	if err == nil {
		t.Fatal("expected command failure")
	}
	if !strings.Contains(err.Error(), failing) {
		t.Fatalf("expected failing command in error, got %v", err)
	}
	if !strings.Contains(err.Error(), "stop") {
		t.Fatalf("expected stderr in error, got %v", err)
	}
}

func TestCheck_CommandSkipsEmptyLines(t *testing.T) {
	t.Parallel()

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND": "  \n\ntrue\n \t\n:",
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_CommandTakesPriorityOverEnvVarCheck(t *testing.T) {
	t.Parallel()

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND":   "true",
		"SEMREL_PLUGIN_ENV_VAR":   "CI",
		"SEMREL_PLUGIN_ENV_VALUE": "true",
		"CI":                      "false",
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_EnvVarMatch(t *testing.T) {
	t.Parallel()

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_ENV_VAR":   "CI",
		"SEMREL_PLUGIN_ENV_VALUE": "true",
		"CI":                      "true",
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_EnvVarMismatch(t *testing.T) {
	t.Parallel()

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_ENV_VAR":   "CI",
		"SEMREL_PLUGIN_ENV_VALUE": "true",
		"CI":                      "false",
	})).Check()
	if err == nil {
		t.Fatal("expected mismatch error")
	}
}

func TestCheck_EmptyValueMatch(t *testing.T) {
	t.Parallel()

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_ENV_VAR":   "OPTIONAL",
		"SEMREL_PLUGIN_ENV_VALUE": "",
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
