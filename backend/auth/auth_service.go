package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
	"valuai/mail"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
)

// issuer is a private constant used as the JWT issuer claim.
const issuer = "valuai.auth"

// AuthenticationConfig contains JWT configuration for authentication.
type AuthenticationConfig struct {
	jwtSecret               string
	tokenExpirationDuration time.Duration
}

// AuthenticationService manages OTP issuance, verification, and token generation.
type AuthenticationService struct {
	emailSender *mail.MailSender
	authConfig  *AuthenticationConfig
	issuedOtps  *cache.Cache
}

// InitAuthConfig initializes an AuthenticationConfig with the given JWT secret and token duration.
func InitAuthConfig(jwtSecret string, tokenExpiration time.Duration) *AuthenticationConfig {
	return &AuthenticationConfig{
		jwtSecret:               jwtSecret,
		tokenExpirationDuration: tokenExpiration,
	}
}

// NewAuthenticationService creates a new Authentication instance with an embedded OTP cache.
func NewAuthenticationService(sender *mail.MailSender, cfg *AuthenticationConfig) *AuthenticationService {
	c := cache.New(10*time.Minute, 1*time.Minute)
	return &AuthenticationService{
		emailSender: sender,
		authConfig:  cfg,
		issuedOtps:  c,
	}
}

// GenerateOtp generates a 6-digit OTP, stores it in cache, and sends it to the user via email.
// The OTP is associated with the provided email and is valid for a short period.
func (auth *AuthenticationService) GenerateOtp(cmd IssueOtpCommand) error {
	otpCode, err := generateSixDigitCode()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	auth.issuedOtps.Set(otpCode, cmd.Email, 10*time.Minute)

	otpTemplate, err := mail.GetTemplate(mail.OtpEmail, cmd.Language)
	if err != nil {
		return err
	}

	emailTitle := otpTemplate.Title
	emailBody, err := otpTemplate.ResolveTemplateContent(otpCode)
	if err != nil {
		return err
	}

	sendEmailCommand := mail.SendMailCommand{
		To:    cmd.Email,
		Title: emailTitle,
		Body:  emailBody,
	}

	if err := auth.emailSender.SendEmail(sendEmailCommand); err != nil {
		return fmt.Errorf("sending OTP to %s failed: %w", cmd.Email, err)
	}

	return nil
}

// VerifyOtp validates the provided OTP against the stored value for the given email.
// If successful, it returns a new access token valid for the configured duration.
func (auth *AuthenticationService) VerifyOtp(verifyOtpCommand VerifyOtpCommand) (*Token, error) {
	value, found := auth.issuedOtps.Get(verifyOtpCommand.Otp)
	if !found {
		return nil, fmt.Errorf("invalid or expired OTP")
	}

	storedEmail, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid cache data for OTP")
	}

	if storedEmail != verifyOtpCommand.Email {
		return nil, fmt.Errorf("OTP does not belong to the provided email")
	}

	auth.issuedOtps.Delete(verifyOtpCommand.Otp)

	var claims = Claims{
		Email:    verifyOtpCommand.Email,
		Language: verifyOtpCommand.Language,
	}

	token, err := auth.generateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return token, nil
}

// generateToken creates a signed JWT access token for the given email with configured expiration.
func (auth *AuthenticationService) generateToken(claims Claims) (*Token, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"iss":      issuer,
		"email":    claims.Email,
		"language": claims.Language,
		"exp":      now.Add(auth.authConfig.tokenExpirationDuration).Unix(),
		"iat":      now.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(auth.authConfig.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken: accessToken,
	}, nil
}

// generateSixDigitCode generates a random 6-digit OTP code as a zero-padded string.
func generateSixDigitCode() (string, error) {
	max := big.NewInt(1000000) // 0 - 999999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
