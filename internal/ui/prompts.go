package ui

import (
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/leeineian/gauth/internal/model"
)

func PromptNewAccount() (*model.Account, error) {
	var (
		issuer  string
		label   string
		secret  string
		otpType string = "totp"
		digits  int    = 6
		algo    string = "sha1"
		period  int    = 30
		counter int    = 0
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Issuer").
				Description("The provider name (e.g. GitHub)").
				Placeholder("GitHub").
				Value(&issuer).
				Validate(required("issuer")),
			huh.NewInput().
				Title("Account ID").
				Description("Username or email associated with the account").
				Placeholder("user@example.com").
				Value(&label).
				Validate(required("account ID")),
			huh.NewInput().
				Title("Secret").
				Description("The Base32 secret key").
				Value(&secret).
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					s = strings.ToUpper(strings.ReplaceAll(s, " ", ""))
					if s == "" {
						return fmt.Errorf("secret is required")
					}
					_, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(s)
					if err != nil {
						// try with padding
						_, err = base32.StdEncoding.DecodeString(s)
					}
					if err != nil {
						return fmt.Errorf("invalid Base32 secret (standard OTP secrets only use A-Z and 2-7)")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Type").
				Options(
					huh.NewOption("TOTP (Time-based)", "totp"),
					huh.NewOption("HOTP (Counter-based)", "hotp"),
				).
				Value(&otpType),
			huh.NewSelect[int]().
				Title("Digits").
				Options(
					huh.NewOption("6 digits", 6),
					huh.NewOption("8 digits", 8),
				).
				Value(&digits),
			huh.NewSelect[string]().
				Title("Algorithm").
				Options(
					huh.NewOption("SHA1", "sha1"),
					huh.NewOption("SHA256", "sha256"),
					huh.NewOption("SHA512", "sha512"),
				).
				Value(&algo),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))

	if otpType == "totp" {
		periodStr := "30"
		err := huh.NewInput().
			Title("Period").
			Description("Time window in seconds (usually 30)").
			Value(&periodStr).
			Validate(func(s string) error {
				var v int
				if _, err := fmt.Sscanf(s, "%d", &v); err != nil || v <= 0 {
					return fmt.Errorf("must be a positive number")
				}
				return nil
			}).Run()
		if err != nil {
			return nil, err
		}
		fmt.Sscanf(periodStr, "%d", &period)
	} else {
		counterStr := "0"
		err := huh.NewInput().
			Title("Initial Counter").
			Description("The starting value for HOTP (usually 0)").
			Value(&counterStr).
			Validate(func(s string) error {
				var v int
				if _, err := fmt.Sscanf(s, "%d", &v); err != nil || v < 0 {
					return fmt.Errorf("must be a non-negative number")
				}
				return nil
			}).Run()
		if err != nil {
			return nil, err
		}
		fmt.Sscanf(counterStr, "%d", &counter)
	}

	return &model.Account{
		Issuer:    issuer,
		Label:     label,
		Secret:    secret,
		Type:      model.OTPType(otpType),
		Digits:    digits,
		Algorithm: algo,
		Period:    int64(period),
		Counter:   int64(counter),
	}, nil
}

func PromptDeleteAccount(accounts []model.Account) (int, error) {
	options := make([]huh.Option[int], 0, len(accounts))
	for i, acc := range accounts {
		options = append(options, huh.NewOption(acc.FullIdentifier(), i))
	}

	var selected int
	err := huh.NewSelect[int]().
		Title("Select account to delete").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return -1, err
	}

	var confirm bool
	err = huh.NewConfirm().
		Title(fmt.Sprintf("Are you sure you want to delete %s?", accounts[selected].FullIdentifier())).
		Value(&confirm).
		Run()

	if err != nil || !confirm {
		return -1, nil
	}

	return selected, nil
}

func PromptPassword(title string) (string, error) {
	var password string
	err := huh.NewInput().
		Title(title).
		EchoMode(huh.EchoModePassword).
		Value(&password).
		Run()
	return password, err
}

func PromptNewPassword() (string, error) {
	var p1, p2 string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("New Master Password").
				Description("Leave empty to remove password protection").
				EchoMode(huh.EchoModePassword).
				Value(&p1).
				Validate(func(s string) error {
					if len(s) > 0 && len(s) < 8 {
						return fmt.Errorf("password must be at least 8 characters")
					}
					return nil
				}),
			huh.NewInput().
				Title("Confirm Password").
				EchoMode(huh.EchoModePassword).
				Value(&p2).
				Validate(func(s string) error {
					if s != p1 {
						return fmt.Errorf("passwords do not match")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	return p1, err
}

func required(name string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%s is required", name)
		}
		return nil
	}
}
