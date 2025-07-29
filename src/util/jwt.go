package utils

import (
	"bookem-user-service/domain"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_PRIVATE_KEY_PATH = os.Getenv("JWT_PRIVATE_KEY_PATH")
var JWT_PUBLIC_KEY_PATH = os.Getenv("JWT_PUBLIC_KEY_PATH")

var CreateJWT = createJWT
var VerifyJWT = verifyJWT

// CreateJWT issues a JSON Web Token with the provided fields.
// The JWT is signed with a private key.
func createJWT(userID int, username string, role domain.UserRole) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"iat":      time.Now().Unix(),
		"user_id":  userID,
		"username": username,
		"role":     role,
	}

	privateKeyData, err := os.ReadFile(JWT_PRIVATE_KEY_PATH)
	if err != nil {
		return "", fmt.Errorf("could not open private key %s: %w", JWT_PRIVATE_KEY_PATH, err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", fmt.Errorf("could not parse private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT checks if the given JWT was provided by this server using a public key.
func verifyJWT(tokenString string) error {
	publicKeyData, err := os.ReadFile(JWT_PUBLIC_KEY_PATH)
	if err != nil {
		return fmt.Errorf("could not open public key %s: %w", JWT_PUBLIC_KEY_PATH, err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return fmt.Errorf("could not parse public key: %w", err)
	}

	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return err
	}
	return nil
}
