package main

import (
	"bytes"
	"testing"
)

func env(kv map[string]string) func(string) string {
	return func(key string) string { return kv[key] }
}

func TestRun_Success(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := run(env(map[string]string{}), &stderr)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("unexpected result: code=%d stderr=%q", code, stderr.String())
	}
}

func TestRun_Failure(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := run(env(map[string]string{
		"SEMREL_PLUGIN_ENV_VAR":   "CI",
		"SEMREL_PLUGIN_ENV_VALUE": "true",
		"CI":                      "false",
	}), &stderr)
	if code != 1 || stderr.Len() == 0 {
		t.Fatalf("unexpected result: code=%d stderr=%q", code, stderr.String())
	}
}
