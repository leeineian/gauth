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

const version = "gauth v1.26.0 (Go 1.25.5)"

var (
	rootCmd = &cobra.Command{
		Use:           "gauth",
		Short:         "gauth is a fast, no-nonsense 2FA for your terminal",
		RunE:          runRoot,
		SilenceUsage:  true,
		SilenceErrors: true,
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
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "watch codes update in real-time")

	var versionFlag, passwdFlag, exportFlag, importFlag, accountFlag, addFlag, deleteFlag bool
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "print version information")
	rootCmd.Flags().BoolVarP(&passwdFlag, "passwd", "p", false, "set or change the master password")
	rootCmd.Flags().BoolVarP(&exportFlag, "export", "e", false, "export accounts to andOTP format")
	rootCmd.Flags().BoolVarP(&importFlag, "import", "i", false, "import accounts from andOTP backups")
	rootCmd.Flags().BoolVarP(&accountFlag, "list", "l", false, "list all accounts")
	rootCmd.Flags().BoolVarP(&addFlag, "add", "a", false, "add a new account")
	rootCmd.Flags().BoolVarP(&deleteFlag, "delete", "d", false, "delete an account")

	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if versionFlag {
			fmt.Println(version)
			os.Exit(0)
		}
		if passwdFlag {
			return passwdCmd.RunE(cmd, args)
		}
		if exportFlag {
			return exportCmd.RunE(cmd, args)
		}
		if importFlag {
			return importCmd.RunE(cmd, args)
		}
		if accountFlag {
			return entryListCmd.RunE(cmd, args)
		}
		if addFlag {
			return entryAddCmd.RunE(cmd, args)
		}
		if deleteFlag {
			return entryDeleteCmd.RunE(cmd, args)
		}
		return nil
	}
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
		fmt.Println("No accounts found. Use 'gauth -a' to add one.")
		return nil
	}

	// Proactively suggest encryption if it's currently plain text
	isEnc, _ := store.IsEncrypted()
	if !isEnc {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("! Your database is currently unencrypted. Run 'gauth -p' to set a master password."))
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
