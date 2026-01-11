package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leeineian/gauth/internal/model"
)

func TestStorage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gauth-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dbFile := filepath.Join(tempDir, "gauth.json")
	s := &Storage{
		baseDir: tempDir,
		dbFile:  dbFile,
	}

	// Test read empty
	accounts, err := s.ReadAccounts("")
	if err != nil {
		t.Errorf("ReadAccounts() error = %v", err)
	}
	if len(accounts) != 0 {
		t.Errorf("expected 0 accounts, got %d", len(accounts))
	}

	// Test write/read with encryption
	pwd := "testpassword"
	testAccounts := []model.Account{
		{
			Issuer:    "TestIssuer",
			Label:     "test@user",
			Secret:    "JBSWY3DPEHPK3PXP",
			Digits:    6,
			Type:      model.TypeTOTP,
			Algorithm: "sha1",
			Period:    30,
		},
	}

	if err := s.WriteAccounts(testAccounts, pwd); err != nil {
		t.Fatalf("WriteAccounts() error = %v", err)
	}

	// Test wrong password
	_, err = s.ReadAccounts("wrong")
	if err == nil {
		t.Error("expected error with wrong password, got nil")
	}

	// Test correct password
	readBack, err := s.ReadAccounts(pwd)
	if err != nil {
		t.Fatalf("ReadAccounts() error = %v", err)
	}

	if len(readBack) != 1 {
		t.Fatalf("expected 1 account, got %d", len(readBack))
	}

	if readBack[0].Issuer != testAccounts[0].Issuer {
		t.Errorf("expected issuer %s, got %s", testAccounts[0].Issuer, readBack[0].Issuer)
	}
}
