package faucet

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTToken(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("This_is_the_secret"))
	assert.NoError(t, err)
	fmt.Println(tokenString)

	// validate the token
	tokenValidation, err := ValidateToken(tokenString, []byte("This_is_the_secret"))
	if err != nil {
		t.Fatalf("unexpected token validation failure - %s", err)
	}
	if !tokenValidation.Valid {
		t.Fatal("token not valid")
	}
}
