// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package plugin

import (
	"os/exec"
	"slices"
	"strconv"
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
	requirePOSIXShell(t)

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND": successfulCommand(),
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_CommandFailure(t *testing.T) {
	t.Parallel()
	requirePOSIXShell(t)

	command := failingCommand("boom", 1)
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
	requirePOSIXShell(t)

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND": successfulCommand() + "\n" + successfulCommand(),
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_MultipleCommandsFailureStopsAtFailingCommand(t *testing.T) {
	t.Parallel()
	requirePOSIXShell(t)

	failing := failingCommand("stop", 2)
	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND": successfulCommand() + "\n" + failing + "\n" + successfulCommand(),
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
	requirePOSIXShell(t)

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND": "  \n\n" + successfulCommand() + "\n \t\n" + successfulCommand(),
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_CommandTakesPriorityOverEnvVarCheck(t *testing.T) {
	t.Parallel()
	requirePOSIXShell(t)

	err := NewWithEnv(env(map[string]string{
		"SEMREL_PLUGIN_COMMAND":   successfulCommand(),
		"SEMREL_PLUGIN_ENV_VAR":   "CI",
		"SEMREL_PLUGIN_ENV_VALUE": "true",
		"CI":                      "false",
	})).Check()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func successfulCommand() string {
	return "true"
}

func failingCommand(message string, exitCode int) string {
	return "printf '" + message + "\\n' >&2; exit " + strconv.Itoa(exitCode)
}

func requirePOSIXShell(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh is unavailable on this host")
	}
}

func TestShellCommandUsesPOSIXShell(t *testing.T) {
	cmd := shellCommand("printf ok")
	if got, want := cmd.Args, []string{"sh", "-c", "printf ok"}; !slices.Equal(got, want) {
		t.Fatalf("shell command args = %q, want %q", got, want)
	}
}

func TestCheckReportsMissingPOSIXShell(t *testing.T) {
	condition := &Condition{
		env: env(map[string]string{"SEMREL_PLUGIN_COMMAND": "true"}),
		shellCommand: func(string) *exec.Cmd {
			return exec.Command("semrel-missing-sh")
		},
	}
	err := condition.Check()
	if err == nil || !strings.Contains(err.Error(), `required shell "sh" is not available on PATH`) {
		t.Fatalf("expected actionable missing sh error, got %v", err)
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
