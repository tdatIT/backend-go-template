package oidc

import (
	"context"
	"log/slog"

	"github.com/tdatIT/backend-go/config"
	"resty.dev/v3"
)

type GoogleIDTokenInfo struct {
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Subject       string `json:"sub"`
	Audience      string `json:"aud"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
}

type GoogleOIDCProvider struct {
	client *resty.Client
}

func NewGoogleOIDCProvider(config *config.ServiceConfig) *GoogleOIDCProvider {
	transportSetting := &resty.TransportSettings{
		DialerTimeout:       0,
		DialerKeepAlive:     0,
		IdleConnTimeout:     0,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 0,
		MaxConnsPerHost:     0,
	}

	newClient := resty.NewWithTransportSettings(transportSetting).
		SetDebug(config.Server.DebugMode).
		SetTimeout(config.Server.HttpClientTimeout)

	return &GoogleOIDCProvider{
		client: newClient,
	}
}

func (p *GoogleOIDCProvider) VerifyIDToken(ctx context.Context, idToken string) (*GoogleIDTokenInfo, error) {
	resp, err := p.client.R().
		SetContext(ctx).
		SetQueryParam("id_token", idToken).
		SetResult(&GoogleIDTokenInfo{}).
		Get("https://oauth2.googleapis.com/tokeninfo")
	if err != nil {
		slog.Error("failed to verify google id token", slog.String("error", err.Error()))
		return nil, err
	}

	if resp.IsError() {
		slog.Error("google tokeninfo endpoint returned error",
			slog.Int("status_code", resp.StatusCode()),
			slog.String("response", resp.String()))
		return nil, err
	}

	info := resp.Result().(*GoogleIDTokenInfo)
	return info, nil
}
