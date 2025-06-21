package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type EncryptedSecret struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type FileStore struct {
	Path string
	mu   sync.Mutex
}

func NewFileStore(path string) *FileStore {
	return &FileStore{Path: path}
}

func (fs *FileStore) loadAll() (map[string]EncryptedSecret, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	secrets := make(map[string]EncryptedSecret)

	file, err := os.ReadFile(fs.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return secrets, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(file, &secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}

func (fs *FileStore) Save(secret EncryptedSecret, overwrite bool) error {
	secrets, err := fs.loadAll()
	if err != nil {
		return err
	}

	if _, exists := secrets[secret.Name]; exists && !overwrite {
		return errors.New("secret with this name already exists (use --force to overwrite)")
	}

	secrets[secret.Name] = secret

	fs.mu.Lock()
	defer fs.mu.Unlock()
	data, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.Path, data, 0600)
}

func (fs *FileStore) Get(name string) (*EncryptedSecret, error) {
	secrets, err := fs.loadAll()
	if err != nil {
		return nil, err
	}

	secret, ok := secrets[name]
	if !ok {
		return nil, errors.New("secret not found")
	}

	return &secret, nil
}
