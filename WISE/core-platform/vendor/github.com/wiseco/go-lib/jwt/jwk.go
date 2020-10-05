package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/wiseco/go-lib/log"
)

// JSON Web Key (JWK): https://tools.ietf.org/html/rfc7517
type Jwk struct {
	KeyID           string `json:"kid"`
	KeyType         string `json:"kty"`
	Algorithm       string `json:"alg"`
	PublicExponent  string `json:"e"`
	PrivateExponent string `json:"d"`
	Modulus         string `json:"n"`
	PublicKeyUse    string `json:"use"`
}

type JwkSet struct {
	Keys []*Jwk `json:"keys"`
}

var (
	clientKeySync  sync.Once
	clientKeys     JwkSet
	supportKeySync sync.Once
	supportKeys    JwkSet
)

func getClientKeys() *JwkSet {
	clientKeySync.Do(func() {
		l := log.NewLogger()

		// Download cognito jwk
		awsRegion := os.Getenv("AWS_REGION")
		userPool := os.Getenv("CLIENT_USER_POOL_ID")
		url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", awsRegion, userPool)

		resp, err := http.Get(url)
		if err != nil {
			l.ErrorD("Error fetching client JWK", log.Fields{"url": url, "error": err})
			return
		}

		defer resp.Body.Close()

		// Decode cognito pool jwk set
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			l.ErrorD("Error fetching client JWK", log.Fields{"error": err})
			return
		}

		err = json.Unmarshal(b, &clientKeys)
		if err != nil {
			l.ErrorD("Error fetching client JWK", log.Fields{"keyset": string(b), "error": err})
			return
		}

		l.InfoD("Client JWK", log.Fields{"keyset": clientKeys})
	})

	return &clientKeys
}

func getClientKey(keyID string) (*Jwk, error) {
	// Lookup by key id
	for _, k := range getClientKeys().Keys {
		if keyID == k.KeyID {
			return k, nil
		}
	}

	log.NewLogger().InfoD("Client key not found", log.Fields{"key": keyID})
	return &Jwk{}, errors.New("key not found")
}

func getSupportKeys() *JwkSet {
	supportKeySync.Do(func() {
		l := log.NewLogger()

		awsRegion := os.Getenv("AWS_REGION")
		userPool := os.Getenv("SUPPORT_USER_POOL_ID")
		url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", awsRegion, userPool)

		resp, err := http.Get(url)
		if err != nil {
			l.ErrorD("Error fetching support JWK", log.Fields{"url": url, "error": err})
			return
		}

		defer resp.Body.Close()

		// Decode cognito pool jwk set
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			l.ErrorD("Error fetching support JWK", log.Fields{"error": err})
			return
		}

		err = json.Unmarshal(b, &supportKeys)
		if err != nil {
			l.ErrorD("Error fetching support JWK", log.Fields{"keyset": string(b), "error": err})
			return
		}

		l.InfoD("Support JWK", log.Fields{"keyset": supportKeys})
	})

	return &supportKeys
}

func getSupportKey(keyID string) (*Jwk, error) {
	// Lookup by key id
	for _, k := range getSupportKeys().Keys {
		if keyID == k.KeyID {
			return k, nil
		}
	}

	log.NewLogger().InfoD("Support key not found", log.Fields{"key": keyID})
	return &Jwk{}, errors.New("key not found")
}
