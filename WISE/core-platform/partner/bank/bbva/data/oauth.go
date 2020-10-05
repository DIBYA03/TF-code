package data

import (
	"os"

	"github.com/wiseco/core-platform/shared"
)

var authSubdomain = map[string]string{
	"sandbox": "sbx-paas",
	"preprod": "pre-paas",
	"prod":    "paas",
}

var oAuthSrv = func() shared.OAuthService {
	// BBVA App Env
	appEnv := os.Getenv("BBVA_APP_ENV")
	if len(appEnv) == 0 {
		appEnv = "sandbox"
	}

	return shared.NewOAuthService(
		shared.OAuthConfig{
			BaseURL:   "https://" + authSubdomain[appEnv] + ".bbvacompass.com/auth/token",
			AppID:     bbvaConfig.appID,
			AppSecret: bbvaConfig.secretOAuthKey,
		},
	)
}()
