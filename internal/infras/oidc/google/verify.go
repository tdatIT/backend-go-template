package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type IDTokenInfo struct {
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Subject       string `json:"sub"`
	Audience      string `json:"aud"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
}

type TokenVerifier struct {
	client *http.Client
}

func NewTokenVerifier() *TokenVerifier {
	return &TokenVerifier{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (v *TokenVerifier) VerifyIDToken(ctx context.Context, idToken string) (*IDTokenInfo, error) {
	if strings.TrimSpace(idToken) == "" {
		return nil, errors.New("id token is required")
	}

	endpoint := "https://oauth2.googleapis.com/tokeninfo"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("id_token", idToken)
	req.URL.RawQuery = q.Encode()

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google tokeninfo error: %s", strings.TrimSpace(string(body)))
	}

	info := new(IDTokenInfo)
	if err := json.Unmarshal(body, info); err != nil {
		return nil, err
	}

	return info, nil
}
