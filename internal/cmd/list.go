package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/leeineian/gauth/internal/storage"
	"github.com/spf13/cobra"
)

var entryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
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
			fmt.Println("No accounts found.")
			return nil
		}

		headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).PaddingRight(1)
		rowStyle := lipgloss.NewStyle().PaddingRight(1)

		tbl := table.New().
			Border(lipgloss.HiddenBorder()).
			Headers("ID", "ISSUER", "LABEL", "TYPE", "DIGITS", "ALGO").
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == 0 {
					return headerStyle
				}
				return rowStyle
			})

		for i, acc := range accounts {
			tbl.Row(
				fmt.Sprintf("%d", i+1),
				acc.Issuer,
				acc.DisplayLabel(),
				strings.ToUpper(string(acc.Type)),
				fmt.Sprintf("%d", acc.Digits),
				strings.ToUpper(acc.Algorithm),
			)
		}

		fmt.Println(tbl.Render())
		return nil
	},
}
