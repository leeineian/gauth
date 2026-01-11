package model

import (
	"fmt"
	"strings"
)

type OTPType string

const (
	TypeTOTP OTPType = "totp"
	TypeHOTP OTPType = "hotp"
)

type Account struct {
	Secret    string                 `json:"secret"`
	Label     string                 `json:"label"`
	Issuer    string                 `json:"issuer"`
	Digits    int                    `json:"digits"`
	Algorithm string                 `json:"algorithm"`
	Counter   int64                  `json:"counter"`
	Period    int64                  `json:"period"`
	Type      OTPType                `json:"type"`
	Misc      map[string]interface{} `json:"misc,omitempty"`
}

func (a *Account) DisplayLabel() string {
	if a.Label != "" && strings.Contains(a.Label, ":") {
		return strings.Split(a.Label, ":")[1]
	}
	if a.Label != "" {
		return a.Label
	}
	return "unnamed"
}

func (a *Account) FullIdentifier() string {
	if a.Issuer != "" {
		return fmt.Sprintf("%s:%s", a.Issuer, a.DisplayLabel())
	}
	return a.DisplayLabel()
}

type OTPResult struct {
	Code      string
	Remaining int64
}

func (a *Account) Validate() error {
	if a.Secret == "" {
		return fmt.Errorf("secret is required")
	}
	if a.Label == "" {
		return fmt.Errorf("label is required")
	}
	if a.Issuer == "" {
		return fmt.Errorf("issuer is required")
	}
	if a.Digits != 6 && a.Digits != 8 {
		return fmt.Errorf("digits must be 6 or 8")
	}
	return nil
}

const (
	DefaultDigits  = 6
	DefaultType    = TypeTOTP
	DefaultCounter = 0
	DefaultAlgo    = "sha1"
	DefaultPeriod  = 30
)
