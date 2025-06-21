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

func TestFileStore_Get(t *testing.T) {
	tmp := "test_secrets.json"
	defer os.Remove(tmp)

	store := NewFileStore(tmp)

	secret := EncryptedSecret{
		Name: "API_KEY",
		Data: []byte("encrypted-api-key"),
	}

	// Save the secret
	if err := store.Save(secret, false); err != nil {
		t.Fatalf("failed to save secret: %v", err)
	}

	// Attempt to retrieve the secret
	retrieved, err := store.Get("API_KEY")
	if err != nil {
		t.Fatalf("failed to get secret: %v", err)
	}

	if retrieved.Name != secret.Name || string(retrieved.Data) != string(secret.Data) {
		t.Errorf("retrieved secret mismatch. got %s, want %s", retrieved.Data, secret.Data)
	}

	// Attempt to retrieve a non-existent secret
	_, err = store.Get("NON_EXISTENT")
	if err == nil {
		t.Error("expected error when getting non-existent secret, got nil")
	}
}

func TestFileStore_Delete(t *testing.T) {
	tmp := "test_secrets.json"
	defer os.Remove(tmp)

	store := NewFileStore(tmp)

	secret := EncryptedSecret{
		Name: "API_KEY",
		Data: []byte("encrypted-api-key"),
	}

	// Save the secret
	if err := store.Save(secret, false); err != nil {
		t.Fatalf("failed to save secret: %v", err)
	}

	// Delete the secret
	if err := store.Delete("API_KEY"); err != nil {
		t.Fatalf("failed to delete secret: %v", err)
	}

	// Attempt to retrieve the deleted secret
	_, err := store.Get("API_KEY")
	if err == nil {
		t.Error("expected error when getting deleted secret, got nil")
	}
}
