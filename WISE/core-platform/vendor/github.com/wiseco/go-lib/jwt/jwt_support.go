package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type SupportAccessTokenClaims struct {
	TokenType  string `json:"token_type"`
	AgentID    string `json:"agent_id"`
	AgentName  string `json:"name"`
	AgentEmail string `json:"email"`
	jwtgo.StandardClaims
}

func ParseSupportAccessToken(accessToken string, verifySig bool) (*jwtgo.Token, error) {
	if verifySig {
		return jwtgo.ParseWithClaims(accessToken, &SupportAccessTokenClaims{}, func(token *jwtgo.Token) (interface{}, error) {
			resp := []byte{}

			switch token.Method {
			case jwtgo.SigningMethodRS256:
				// Cognito
				keyID, ok := token.Header["kid"].(string)
				if !ok {
					return resp, errors.New("key id missing")
				}

				// Find key by id
				k, err := getSupportKey(keyID)
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
		token, _, err := parser.ParseUnverified(accessToken, &SupportAccessTokenClaims{})
		return token, err
	}
}
