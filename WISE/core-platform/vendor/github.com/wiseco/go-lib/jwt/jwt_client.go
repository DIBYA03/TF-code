package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type ClientAccessTokenClaims struct {
	ConsumerID string `json:"consumer_id"`
	TokenType  string `json:"token_type"`
	jwtgo.StandardClaims

	// User id not in subject for cognito
	UserID string `json:"user_id"`
}

const (
	// 15 minute expiration window
	clientAccessTokenExpWindow = 900
)

// TODO: Extend to pass in signing config including signing algo and signing key
func GenerateClientAccessToken(accessTokenID, userID, consumerID, issuer string) (*jwtgo.Token, string, error) {
	issued := time.Now().UTC().Unix()
	claims := &ClientAccessTokenClaims{
		ConsumerID: consumerID,
		TokenType:  TokenTypeClientAccess,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: issued + clientAccessTokenExpWindow,
			Id:        accessTokenID,
			IssuedAt:  issued,
			Issuer:    issuer,
			Subject:   userID,
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)

	signingKey := os.Getenv("JWT_SIGNING_KEY_HS256")
	if signingKey == "" {
		return token, "", errors.New("missing signing key")
	}

	b, err := base64.StdEncoding.DecodeString(signingKey)
	if err != nil {
		return token, "", err
	}

	tokenString, err := token.SignedString(b)
	return token, tokenString, err
}

func ParseClientAccessToken(accessToken string, verifySig bool) (*jwtgo.Token, error) {
	if verifySig {
		return jwtgo.ParseWithClaims(accessToken, &ClientAccessTokenClaims{}, func(token *jwtgo.Token) (interface{}, error) {
			resp := []byte{}

			switch token.Method {
			case jwtgo.SigningMethodHS256:
				// Default signing method
				signingKey := os.Getenv("JWT_SIGNING_KEY_HS256")
				if signingKey == "" {
					return resp, errors.New("missing signing key")
				}

				b, err := base64.StdEncoding.DecodeString(signingKey)
				if err != nil {
					return resp, err
				}

				return b, nil
			case jwtgo.SigningMethodRS256:
				// Cognito
				keyID, ok := token.Header["kid"].(string)
				if !ok {
					return resp, errors.New("key id missing")
				}

				// Find key by id
				k, err := getClientKey(keyID)
				if err != nil {
					return resp, err
				}

				// Decode modulus
				b, err := base64.RawURLEncoding.DecodeString(k.Modulus)
				if err != nil {
					return resp, err
				}

				// Base64: AQAB
				// Hexadecimal: 01 00 01
				// Decimal: 65537
				if k.PublicExponent != "AQAB" {
					return resp, errors.New("invalid public exponent 'e'")
				}

				pk := &rsa.PublicKey{
					N: new(big.Int).SetBytes(b),
					E: 65537,
				}

				return pk, nil
			default:
				return resp, fmt.Errorf("unsupported signing method: %v", token.Header["alg"])
			}
		})
	} else {
		parser := &jwtgo.Parser{}
		token, _, err := parser.ParseUnverified(accessToken, &ClientAccessTokenClaims{})
		return token, err
	}
}

type ClientRefreshTokenClaims struct {
	ClientKey  string `json:"client_key"`
	ConsumerID string `json:"consumer_id"`
	TokenType  string `json:"token_type"`
	jwtgo.StandardClaims
}

const (
	// 6 month expiration window
	refreshTokenExpWindow = 15552000
)

func GenerateClientRefreshToken(refreshTokenID, clientKey, userID, consumerID, issuer string) (*jwtgo.Token, string, error) {
	issued := time.Now().UTC().Unix()
	claims := &ClientRefreshTokenClaims{
		ClientKey:  clientKey,
		ConsumerID: consumerID,
		TokenType:  TokenTypeClientRefresh,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: issued + refreshTokenExpWindow,
			Id:        refreshTokenID,
			IssuedAt:  issued,
			Issuer:    issuer,
			Subject:   userID,
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	signingKey := os.Getenv("JWT_SIGNING_KEY_HS256")
	if signingKey == "" {
		return token, "", errors.New("missing signing key")
	}

	b, err := base64.StdEncoding.DecodeString(signingKey)
	if err != nil {
		return token, "", err
	}

	tokenString, err := token.SignedString(b)
	return token, tokenString, err
}

func ParseClientRefreshToken(refreshToken string) (*jwtgo.Token, error) {
	return jwtgo.ParseWithClaims(refreshToken, &ClientRefreshTokenClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		method, ok := token.Method.(*jwtgo.SigningMethodHMAC)
		if !ok {
			return []byte{}, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		switch method {
		case jwtgo.SigningMethodHS256:
			break
		default:
			return []byte{}, fmt.Errorf("unsupported signing method: %v", token.Header["alg"])
		}

		signingKey := os.Getenv("JWT_SIGNING_KEY_HS256")
		if signingKey == "" {
			return []byte{}, errors.New("missing signing key")
		}

		b, err := base64.StdEncoding.DecodeString(signingKey)
		if err != nil {
			return "", err
		}

		return b, nil
	})
}
