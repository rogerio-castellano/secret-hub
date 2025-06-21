package storage

import (
	"os"
	"testing"
)

func TestFileStore_SaveAndLoad(t *testing.T) {
	tmp := "test_secrets.json"
	defer os.Remove(tmp)

	store := NewFileStore(tmp)

	secret := EncryptedSecret{
		Name: "API_KEY",
		Data: []byte("encrypted-api-key"),
	}

	// Save
	err := store.Save(secret, false)
	if err != nil {
		t.Fatalf("failed to save secret: %v", err)
	}

	// Save again without force
	err = store.Save(secret, false)
	if err == nil {
		t.Error("expected error on duplicate secret without --force, got nil")
	}

	// Save again with force (should pass)
	err = store.Save(secret, true)
	if err != nil {
		t.Fatalf("failed to overwrite secret with --force: %v", err)
	}

	//Load secrets
	all, err := store.loadAll()
	if err != nil {
		t.Fatalf("failed to load secrets: %v", err)
	}

	got, exists := all["API_KEY"]
	if !exists {
		t.Error("secret not found after storing")
	}

	if string(got.Data) != string(secret.Data) {
		t.Errorf("stored data mismatch. got %s, want %s", secret.Data, got.Data)
	}
}
