package cmd

import (
	"fmt"

	"github.com/leeineian/gauth/internal/storage"
	"github.com/leeineian/gauth/internal/ui"
	"github.com/spf13/cobra"
)

var masterPassword string

func getOrPromptPassword(store *storage.Storage) (string, error) {
	if masterPassword != "" {
		return masterPassword, nil
	}

	isEnc, err := store.IsEncrypted()
	if err != nil {
		return "", err
	}

	if !isEnc {
		return "", nil // Plain text
	}

	pwd, err := ui.PromptPassword("Enter Master Password")
	if err != nil {
		return "", err
	}
	masterPassword = pwd
	return pwd, nil
}

var passwdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "Set or change the master password",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStorage()
		if err != nil {
			return err
		}

		currentPwd, err := getOrPromptPassword(store)
		if err != nil {
			return err
		}

		accounts, err := store.ReadAccounts(currentPwd)
		if err != nil {
			return err
		}

		newPwd, err := ui.PromptNewPassword()
		if err != nil {
			return err
		}

		if err := store.WriteAccounts(accounts, newPwd); err != nil {
			return err
		}

		masterPassword = newPwd
		if newPwd == "" {
			fmt.Println("✓ Master password removed. Database is now unencrypted.")
		} else {
			fmt.Println("✓ Master password updated successfully!")
		}
		return nil
	},
}
