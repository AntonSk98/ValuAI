package auth

import "valuai/common"

// Token represents an access token issued after successful OTP verification.
type Token struct {
	AccessToken string `json:"access_token"`
}

// Claims holds the information embedded in the JWT access token.
type Claims struct {
	Email    string
	Language common.Language
}
