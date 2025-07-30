package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestSecretHubIntegration(t *testing.T) {
	// Setup temp files
	tmpDir := t.TempDir()
	keyFile := tmpDir + "/key.bin"
	secretName := "INTEGRATION_SECRET"
	secretValue := "IntegrationTest123!"
	storeFile := tmpDir + "/store.json"

	// 1. Generate key
	err := runCommand("go", "run", ".", "generate-key", "--out", keyFile)
	if err != nil {
		t.Fatalf("generate-key failed: %v", err)
	}

	// 2. Store secret
	err = runCommand("go", "run", ".", "store",
		"--name", secretName,
		"--value", secretValue,
		"--key", keyFile,
		"--storage", storeFile)
	if err != nil {
		t.Fatalf("store failed: %v", err)
	}

	// 3. List secrets
	output, err := runCommandOutput("go", "run", ".", "list", "--storage", storeFile)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(output, secretName) {
		t.Errorf("list output missing secret: %s", output)
	}

	// 4. Get secret
	output, err = runCommandOutput("go", "run", ".", "get",
		"--name", secretName,
		"--key", keyFile,
		"--storage", storeFile)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if strings.TrimSpace(output) != secretValue {
		t.Errorf("decrypted value mismatch: got %s, want %s", output, secretValue)
	}

	// 5. Delete secret
	err = runCommand("go", "run", ".", "delete",
		"--name", secretName,
		"--storage", storeFile)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// 6. Ensure secret is deleted
	output, err = runCommandOutput("go", "run", ".", "list", "--storage", storeFile)
	if err != nil {
		t.Fatalf("list after delete failed: %v", err)
	}
	if strings.Contains(output, secretName) {
		t.Errorf("secret still present after delete: %s", output)
	}
}

// Helpers

func runCommand(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommandOutput(args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}
