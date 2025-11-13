package auth

// Token represents an access token issued after successful OTP verification.
type Token struct {
	AccessToken string `json:"access_token"`
}
