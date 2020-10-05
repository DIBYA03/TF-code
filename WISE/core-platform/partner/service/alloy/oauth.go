/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package alloy

import (
	"os"

	"github.com/wiseco/core-platform/shared"
)

var authSubdomain = map[string]string{
	"sandbox": "sandbox",
	"prod":    "api",
}

var oAuthSrv = func() shared.OAuthService {
	// BBVA App Env
	appEnv := os.Getenv("ALLOY_APP_ENV")
	if len(appEnv) == 0 {
		appEnv = "sandbox"
	}

	return shared.NewOAuthService(
		shared.OAuthConfig{
			BaseURL:   "https://" + authSubdomain[appEnv] + ".alloy.co/v1/oauth/bearer",
			AppID:     alloyConfig.appID,
			AppSecret: alloyConfig.secretOAuthKey,
		},
	)
}()
