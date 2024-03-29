package webtoken

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// ErrNotValidYet is when a token is used prior to the issue date.
	ErrNotValidYet = errors.New("token is not valid yet")
	// ErrExpired is when a token is used after the expiration date.
	ErrExpired = errors.New("token is expired")
	// ErrMalformed is when a token is malformed.
	ErrMalformed = errors.New("token is malformed")
	// ErrSignatureInvalid is when a signature is invalid.
	ErrSignatureInvalid = errors.New("signature is invalid")
	// ErrAudienceInvalid is when the audience is invalid.
	ErrAudienceInvalid = errors.New("audience is invalid")
	// ErrExpirationInvalid is when the expiration is invalid.
	ErrExpirationInvalid = errors.New("expiration is invalid")
	// ErrIssuedAtInvalid is when the issued date is invalid.
	ErrIssuedAtInvalid = errors.New("issue date is invalid")
	// ErrNotBeforeInvalid is when the before date is invalid.
	ErrNotBeforeInvalid = errors.New("before date is invalid")
	// ErrSecretTooShort is when the secret is not long enough.
	ErrSecretTooShort = errors.New("secret must be 256 bit (32 bytes)")
)

// Configuration contains the JWT dependencies.
type Configuration struct {
	clock   IClock
	secret  []byte
	timeout time.Duration
}

// New creates a new JWT configuration.
func New(secret []byte, sessionTimeout time.Duration) *Configuration {
	return &Configuration{
		clock:   new(Clock),
		secret:  secret,
		timeout: sessionTimeout,
	}
}

// SetClock will set the clock.
func (c *Configuration) SetClock(clock IClock) {
	c.clock = clock
}

// randomID generates a UUID for use as an ID.
func randomID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// Generate will generate a JWT.
func (c *Configuration) Generate(userID string) (string, error) {
	// Ensure a secret is present.
	if len(c.secret) < 32 {
		return "", ErrSecretTooShort
	}

	// Get the current time.
	now := c.clock.Now()

	// Generate a unique ID.
	unique, err := randomID()
	if err != nil {
		return "", err
	}

	// Create the claims.
	claims := &jwt.StandardClaims{
		Id:        unique,
		Audience:  userID,
		NotBefore: now.Unix(),
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(c.timeout).Unix(),
	}

	// Create the token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signe the token.
	return token.SignedString([]byte(c.secret))
}

// Verify will ensure a JWT is valid and returns the audience if successful.
func (c *Configuration) Verify(s string) (string, error) {
	token, err := jwt.ParseWithClaims(s, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.secret), nil
	})

	if err == nil {
		// If a token is valid, return the audience.
		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			if claims.ExpiresAt == 0 {
				return "", ErrExpirationInvalid
			} else if claims.NotBefore == 0 {
				return "", ErrNotBeforeInvalid
			} else if claims.IssuedAt == 0 {
				return "", ErrIssuedAtInvalid
			} else if len(claims.Audience) == 0 {
				return "", ErrAudienceInvalid
			}
			return claims.Audience, nil
		}
	}

	// Handle the error.
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "", ErrMalformed
		} else if ve.Errors&(jwt.ValidationErrorSignatureInvalid) != 0 {
			return "", ErrSignatureInvalid
		} else if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			return "", ErrExpired
		} else if ve.Errors&(jwt.ValidationErrorNotValidYet) != 0 {
			return "", ErrNotValidYet
		}
	}

	return "", err
}
