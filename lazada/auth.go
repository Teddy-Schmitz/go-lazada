package lazada

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// The Auth Service deals with doing the OAuth flow and exchanging codes for tokens.
// It also lets you refresh tokens in order to get new credentials.
type AuthService service

// Token is the data returned when doing an Oauth Flow through the open platform
type Token struct {
	AccountID       string `json:"account_id"`
	Account         string `json:"account"`
	Country         string `json:"country"`
	AccountPlatform string `json:"account_platform"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`

	ExpiresIn        int `json:"expires_in"`
	RefreshExpiresIn int `json:"refresh_expires_in"`

	retrievedAt time.Time
}

// ExpiresAt tells you the point in time when this token will expire
func (t *Token) ExpiresAt() time.Time {
	return t.calculateExpires(t.ExpiresIn)
}

// RefreshExpiresAt tells you the point in time when the refresh token will expire
func (t *Token) RefreshExpiresAt() time.Time {
	return t.calculateExpires(t.RefreshExpiresIn)
}

// Valid tells you if the token is valid right now
func (t *Token) Valid() bool {
	if t.AccessToken == "" || time.Now().After(t.ExpiresAt()) {
		return false
	}
	return true
}

func (t *Token) calculateExpires(exp int) time.Time {
	return t.retrievedAt.Add(time.Second * time.Duration(exp))
}

// AuthURL returns the URL you should use to start the OAuth flow
// It takes in the URL that the user should be returned to as redirect
// and a state variable which should be a random string
func (a *AuthService) AuthURL(redirect, state string) string {
	baseURL, _ := url.Parse("https://auth.lazada.com/oauth/authorize")

	q := baseURL.Query()
	q.Set("client_id", a.client.appKey)
	q.Set("redirect_uri", redirect)
	q.Set("response_type", "code")
	q.Set("state", state)

	baseURL.RawQuery = q.Encode()
	return baseURL.String()
}

// Exchange sends the received oauth code to the open platform and returns a token
func (a *AuthService) Exchange(ctx context.Context, code string) (*Token, error) {
	req, err := a.client.NewRequest("GET",
		fmt.Sprintf("%s?code=%s", apiNames["AccessToken"], code), nil)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = a.client.Do(ctx, req, &buf)
	if err != nil {
		return nil, err
	}

	t := &Token{retrievedAt: time.Now()}
	if err := json.NewDecoder(&buf).Decode(t); err != nil {
		return nil, errors.Wrap(err, "cant unmarshal token")
	}

	return t, nil
}

// Refresh sends the refresh token and returns the new refreshed token
func (a *AuthService) Refresh(ctx context.Context, token string) (*Token, error) {
	req, err := a.client.NewRequest("GET",
		fmt.Sprintf("%s?refresh_token=%s", apiNames["RefreshToken"], token), nil)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = a.client.Do(ctx, req, &buf)
	if err != nil {
		return nil, err
	}

	t := &Token{retrievedAt: time.Now()}
	if err := json.NewDecoder(&buf).Decode(t); err != nil {
		return nil, errors.Wrap(err, "cant unmarshal token")
	}

	return t, nil
}
