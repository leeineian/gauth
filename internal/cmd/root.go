package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/leeineian/gauth/internal/model"
	"github.com/leeineian/gauth/internal/service"
	"github.com/leeineian/gauth/internal/storage"
	"github.com/leeineian/gauth/internal/ui"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gauth",
		Short: "gauth is a simple CLI for generating 2FA codes",
		Example: `  gauth             # Show all codes
  gauth entry add   # Add a new account
  gauth entry list  # Manage existing accounts`,
		RunE:          runRoot,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	entryCmd = &cobra.Command{
		Use:   "entry",
		Short: "Manage account entries",
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gauth v1.0.0 (Go 1.25.5)")
		},
	}

	watchFlag bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true).
			MarginTop(1)

		fmt.Fprintln(os.Stderr, errorStyle.Render("Error:"), err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(entryCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(passwdCmd)

	entryCmd.AddCommand(entryListCmd)
	entryCmd.AddCommand(entryAddCmd)
	entryCmd.AddCommand(entryDeleteCmd)

	rootCmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "watch codes update in real-time")
}

func runRoot(cmd *cobra.Command, args []string) error {
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
		fmt.Println("No accounts found. Use 'gauth entry add' to add one.")
		return nil
	}

	// Proactively suggest encryption if it's currently plain text
	isEnc, _ := store.IsEncrypted()
	if !isEnc {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("! Your database is currently unencrypted. Run 'gauth passwd' to set a master password."))
	}

	if watchFlag {
		return ui.RunLiveView(accounts)
	}

	otpSvc := service.NewOTPService()

	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).PaddingRight(1)
	rowStyle := lipgloss.NewStyle().PaddingRight(1)

	tbl := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers("ISSUER", "LABEL", "TYPE", "CODE", "REMAINING").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			return rowStyle
		})

	for _, acc := range accounts {
		res, err := otpSvc.Generate(&acc)
		code := "ERROR"
		remaining := "-"

		if err == nil {
			codeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
			if res.Remaining < 5 && acc.Type == model.TypeTOTP {
				codeStyle = codeStyle.Foreground(lipgloss.Color("9"))
			}
			code = codeStyle.Render(res.Code)
			if acc.Type == model.TypeTOTP {
				remaining = fmt.Sprintf("%ds", res.Remaining)
			}
		}

		tbl.Row(
			acc.Issuer,
			acc.DisplayLabel(),
			strings.ToUpper(string(acc.Type)),
			code,
			remaining,
		)
	}

	fmt.Println(tbl.Render())
	return nil
}
