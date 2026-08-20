package config

import (
	"os"
	"testing"
)

// clearEnv ensures a clean slate between test cases so leftover env vars
// from one test (or the developer's shell) don't leak into another.
func clearEnv(t *testing.T) {
	t.Helper()
	keys := []string{"BACKEND_WS_PORT", "SQLITE_DB_PATH", "TEMP_IMAGE_DIR", "AI_MODELS_PATH"}
	for _, k := range keys {
		original, existed := os.LookupEnv(k)
		os.Unsetenv(k)
		t.Cleanup(func() {
			if existed {
				os.Setenv(k, original)
			}
		})
	}
}

func TestLoadReturnsDefaultsWhenEnvVarsAreUnset(t *testing.T) {
	clearEnv(t)

	cfg := Load()

	if cfg.WSPort != "8080" {
		t.Errorf("expected default WSPort '8080', got %q", cfg.WSPort)
	}
	if cfg.DBPath != "./foto-facil.db" {
		t.Errorf("expected default DBPath './foto-facil.db', got %q", cfg.DBPath)
	}
	if cfg.TempImageDir != "./tmp_processed_images" {
		t.Errorf("expected default TempImageDir './tmp_processed_images', got %q", cfg.TempImageDir)
	}
	if cfg.AIModelsPath != "./backend/models" {
		t.Errorf("expected default AIModelsPath './backend/models', got %q", cfg.AIModelsPath)
	}
}

func TestLoadHonorsExistingEnvironmentVariables(t *testing.T) {
	clearEnv(t)
	os.Setenv("BACKEND_WS_PORT", "9090")
	os.Setenv("SQLITE_DB_PATH", "/tmp/custom.db")

	cfg := Load()

	if cfg.WSPort != "9090" {
		t.Errorf("expected overridden WSPort '9090', got %q", cfg.WSPort)
	}
	if cfg.DBPath != "/tmp/custom.db" {
		t.Errorf("expected overridden DBPath '/tmp/custom.db', got %q", cfg.DBPath)
	}
}
