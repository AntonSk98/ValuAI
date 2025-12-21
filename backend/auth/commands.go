package auth

import (
	"errors"
	"valuai/common"
)

// IssueOtpCommand contains the email and language for OTP delivery.
type IssueOtpCommand struct {
	Email    string          `json:"email"`
	Language common.Language `json:"language"`
}

func (cmd IssueOtpCommand) Validate() error {
	if cmd.Email == "" {
		return errors.New("email is a required parameter")
	}

	if !common.IsLanguageSupported(cmd.Language) {
		return errors.New("unsupported language")
	}

	return nil
}

// VerifyOtpCommand contains the email and OTP code for verification.
type VerifyOtpCommand struct {
	Email    string          `json:"email"`
	Otp      string          `json:"otp"`
	Language common.Language `json:"language"`
}

func (cmd VerifyOtpCommand) Validate() error {
	if cmd.Email == "" {
		return errors.New("email is a required parameter")
	}

	if cmd.Otp == "" {
		return errors.New("otp is a required parameter")
	}

	if !common.IsLanguageSupported(cmd.Language) {
		return errors.New("unsupported language")
	}

	return nil
}
