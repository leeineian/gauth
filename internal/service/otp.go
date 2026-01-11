package service

import (
	"strings"
	"time"

	"github.com/grijul/otpgen"
	"github.com/leeineian/gauth/internal/model"
)

type OTPService struct{}

func NewOTPService() *OTPService {
	return &OTPService{}
}

func (s *OTPService) Generate(acc *model.Account) (*model.OTPResult, error) {
	if strings.ToLower(string(acc.Type)) == "hotp" {
		return s.generateHOTP(acc)
	}
	return s.generateTOTP(acc)
}

func (s *OTPService) generateTOTP(acc *model.Account) (*model.OTPResult, error) {
	now := time.Now().Unix()
	period := acc.Period
	if period == 0 {
		period = 30
	}

	t := &otpgen.TOTP{
		Secret:    acc.Secret,
		Digits:    acc.Digits,
		Algorithm: acc.Algorithm,
		Period:    period,
		UnixTime:  now,
	}

	code, err := t.Generate()
	if err != nil {
		return nil, err
	}

	return &model.OTPResult{
		Code:      code,
		Remaining: period - (now % period),
	}, nil
}

func (s *OTPService) generateHOTP(acc *model.Account) (*model.OTPResult, error) {
	h := &otpgen.HOTP{
		Secret:  acc.Secret,
		Digits:  acc.Digits,
		Counter: acc.Counter,
	}

	code, err := h.Generate()
	if err != nil {
		return nil, err
	}

	return &model.OTPResult{
		Code:      code,
		Remaining: 0,
	}, nil
}
