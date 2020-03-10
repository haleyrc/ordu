package ordu_test

import (
	"os"
	"testing"

	"github.com/haleyrc/ordu"
)

func TestLoadEnvironment(t *testing.T) {
	if err := os.Setenv("test_key", "test_val"); err != nil {
		t.Fatal(err)
	}
	env, err := ordu.LoadEnvironmentWithDefaults(ordu.Environment{
		"test_key":  "defaut_value",
		"test_key2": "default_value",
	})
	if err != nil {
		t.Fatal(err)
	}
	{
		got := env.Get("test_key")
		if got != "test_val" {
			t.Errorf("env.Get(%q) = %q, expected %q", "test_key", got, "test_val")
		}
	}
	{
		got := env.Get("test_key2")
		if got != "default_value" {
			t.Errorf("env.Get(%q) = %q, expected %q", "test_key2", got, "default_value")

		}
	}
	{
		got := env.Get("test_key3")
		if got != "" {
			t.Errorf("env.Get(%q) = %q, expected %q", "test_key", got, "")
		}
	}
}
