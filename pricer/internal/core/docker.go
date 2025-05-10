package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultSecretPath = "/run/secrets"

func ReadSecret(key string) (string, error) {
	secretPath := filepath.Join(defaultSecretPath, key)

	data, err := os.ReadFile(secretPath)
	if err != nil {
		return "", fmt.Errorf("failed to read secret %s: %w", key, err)
	}

	return strings.TrimSpace(string(data)), nil
}

func MustReadSecret(key string) string {
	val, err := ReadSecret(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "critical secret %q could not be loaded: %v\n", key, err)
		os.Exit(1)
	}
	return val
}

