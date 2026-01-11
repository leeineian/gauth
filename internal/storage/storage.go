package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/leeineian/gauth/internal/model"
)

type Storage struct {
	baseDir string
	dbFile  string
}

func NewStorage() (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}

	baseDir := filepath.Join(home, ".gauth")
	dbFile := filepath.Join(baseDir, "gauth.json")

	return &Storage{
		baseDir: baseDir,
		dbFile:  dbFile,
	}, nil
}

func (s *Storage) EnsureDir() error {
	return os.MkdirAll(s.baseDir, 0700)
}

func (s *Storage) IsEncrypted() (bool, error) {
	data, err := os.ReadFile(s.dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	var accounts []model.Account
	err = json.Unmarshal(data, &accounts)
	return err != nil, nil // If error, it's likely encrypted
}

func (s *Storage) ReadAccounts(password string) ([]model.Account, error) {
	data, err := os.ReadFile(s.dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Account{}, nil
		}
		return nil, fmt.Errorf("failed to read database: %w", err)
	}

	var accounts []model.Account
	if err := json.Unmarshal(data, &accounts); err == nil {
		return accounts, nil
	}

	if password == "" {
		return nil, fmt.Errorf("database is encrypted, please provide a password")
	}

	decrypted, err := decrypt(data, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt database (wrong password?): %w", err)
	}

	if err := json.Unmarshal(decrypted, &accounts); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted database: %w", err)
	}

	return accounts, nil
}

func (s *Storage) WriteAccounts(accounts []model.Account, password string) error {
	if err := s.EnsureDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode accounts: %w", err)
	}

	finalData := data
	if password != "" {
		encrypted, err := encrypt(data, password)
		if err != nil {
			return fmt.Errorf("failed to encrypt accounts: %w", err)
		}
		finalData = encrypted
	}

	// Atomic write using temp file rename
	tmpFile := s.dbFile + ".tmp"
	if err := os.WriteFile(tmpFile, finalData, 0600); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tmpFile, s.dbFile); err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}

func (s *Storage) GetFileLocation() string {
	return s.dbFile
}
