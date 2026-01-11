package cmd

import (
	"fmt"

	"github.com/leeineian/gauth/internal/storage"
	"github.com/leeineian/gauth/internal/ui"
	"github.com/spf13/cobra"
)

var entryAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new account interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		acc, err := ui.PromptNewAccount()
		if err != nil {
			return err
		}

		store, err := storage.NewStorage()
		if err != nil {
			return err
		}

		pwd, err := getOrPromptPassword(store)
		if err != nil {
			return err
		}

		accounts, err := store.ReadAccounts(pwd)
		if err != nil {
			return err
		}

		for _, a := range accounts {
			if a.FullIdentifier() == acc.FullIdentifier() {
				return fmt.Errorf("account already exists: %s", acc.FullIdentifier())
			}
		}

		accounts = append(accounts, *acc)
		if err := store.WriteAccounts(accounts, pwd); err != nil {
			return err
		}

		fmt.Printf("\n✓ Account for %s added successfully!\n", acc.FullIdentifier())
		return nil
	},
}

var entryDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete an account",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStorage()
		if err != nil {
			return err
		}

		pwd, err := getOrPromptPassword(store)
		if err != nil {
			return err
		}

		accounts, err := store.ReadAccounts(pwd)
		if err != nil {
			return err
		}

		if len(accounts) == 0 {
			fmt.Println("No accounts to delete.")
			return nil
		}

		selectedIdx, err := ui.PromptDeleteAccount(accounts)
		if err != nil {
			return err
		}

		if selectedIdx < 0 || selectedIdx >= len(accounts) {
			return nil // Cancelled
		}

		deleted := accounts[selectedIdx]
		accounts = append(accounts[:selectedIdx], accounts[selectedIdx+1:]...)

		if err := store.WriteAccounts(accounts, pwd); err != nil {
			return err
		}

		fmt.Printf("✓ Deleted account: %s\n", deleted.FullIdentifier())
		return nil
	},
}
