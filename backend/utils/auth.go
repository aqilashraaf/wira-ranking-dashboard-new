package utils

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret    = []byte("your-secret-key") // TODO: Move to environment variable
	tokenExpiry  = time.Hour * 24            // 24 hours
)

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(userID int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(tokenExpiry).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// HashPassword creates a bcrypt hash of a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate2FASecret generates a new TOTP secret
func Generate2FASecret() (string, string, error) {
	// Generate random bytes for the secret
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", "", err
	}

	// Encode the secret in base32
	secret := base32.StdEncoding.EncodeToString(bytes)

	// Generate the QR code URL
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "WIRA Dashboard",
		AccountName: "user@example.com",
		SecretSize:  20,
	})
	if err != nil {
		return "", "", err
	}

	return secret, key.URL(), nil
}

// Validate2FACode validates a TOTP code
func Validate2FACode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateRefreshToken generates a new refresh token
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(bytes), nil
}
