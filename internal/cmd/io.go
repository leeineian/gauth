package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/leeineian/gauth/internal/provider/andotp"
	"github.com/leeineian/gauth/internal/storage"
	"github.com/leeineian/gauth/internal/ui"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import accounts from andOTP backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")

		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		var password string
		// Check if the backup is encrypted (invalid JSON indicates encryption)
		var test interface{}
		if err := json.Unmarshal(data, &test); err != nil {
			p, err := ui.PromptPassword("Enter Backup Decryption Password")
			if err != nil {
				return err
			}
			password = p
		}

		prov := andotp.New()
		accounts, err := prov.Import(filePath, password)
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

		existing, err := store.ReadAccounts(pwd)
		if err != nil {
			return err
		}

		// Simple duplicate check
		seen := make(map[string]bool)
		for _, e := range existing {
			seen[e.FullIdentifier()] = true
		}

		newCount := 0
		for _, a := range accounts {
			if !seen[a.FullIdentifier()] {
				existing = append(existing, a)
				seen[a.FullIdentifier()] = true
				newCount++
			}
		}

		if newCount == 0 {
			fmt.Println("No new accounts found in backup (all already exist).")
			return nil
		}

		if err := store.WriteAccounts(existing, pwd); err != nil {
			return err
		}

		fmt.Printf("✓ Imported %d new accounts (skipped %d duplicates)\n", newCount, len(accounts)-newCount)
		return nil
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export accounts to andOTP format",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")

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

		password, err := ui.PromptPassword("Enter Password to Encrypt Export (leave empty for plain text)")
		if err != nil {
			return err
		}

		prov := andotp.New()
		data, err := prov.Export(accounts, password)
		if err != nil {
			return err
		}

		if err := os.WriteFile(filePath, data, 0600); err != nil {
			return err
		}

		fmt.Printf("✓ Exported %d accounts to %s\n", len(accounts), filePath)
		return nil
	},
}

func init() {
	importCmd.Flags().StringP("file", "f", "", "Backup file to import")
	importCmd.MarkFlagRequired("file")

	exportCmd.Flags().StringP("file", "f", "gauth_backup.json", "Output file")
}
