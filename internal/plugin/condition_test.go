package plugin

import "testing"

func env(kv map[string]string) func(string) string {
	return func(key string) string { return kv[key] }
}

func TestCheck_PermissiveWithoutEnvVarName(t *testing.T) {
	t.Parallel()

	if err := NewWithEnv(env(map[string]string{})).Check(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_Match(t *testing.T) {
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

func TestCheck_Mismatch(t *testing.T) {
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
