package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/calamity-of-subterfuge/cos/pkg/utils"
	"golang.org/x/net/http2"
)

type loginResult struct {
	Token     string  `json:"token"`
	ExpiresAt float64 `json:"expires_at"`
}

// AuthToken is the result of successfully logging in. It provides an arbitrary
// token and the time at which the token will stop working if the user does not
// logout sooner.
type AuthToken struct {
	// Token is an arbitrary secret that is provided by the calamity of subterfuge
	// website which we can pass to future requests as the bearer token. The format
	// of this token, if it's formatted at all, is not guarranteed. Hence this MUST
	// be treated as an opaque string.
	Token string

	// ExpiresAt is the time at which this AuthToken will stop working if no other
	// actions are taken. Note that this is just a hint from the server; the server
	// may decide to invalidate tokens at any time for any reason.
	ExpiresAt time.Time
}

// Login will login to the calamity of subterfuge account with the given email,
// using the given grant and secret. Note that this is not the same technique as
// humans use during the website, although it is the same endpoint. Any number
// of grant identifier and secret pairs can be created for an account after you
// signup via the website.
func Login(email, grantIden, secret string) (*AuthToken, error) {
	body := map[string]interface{}{
		"email":      email,
		"password":   secret,
		"grant_iden": grantIden,
	}
	bodyMarshalled, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	client := &http.Client{
		Transport: &http2.Transport{},
	}
	var resp *http.Response
	resp, err = client.Post(
		API_BASE+"/api/1/auth/sessions",
		"application/json",
		bytes.NewBuffer(bodyMarshalled),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to POST: %w", err)
	}

	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close body: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d (body: %s)", resp.StatusCode, string(respBody))
	}

	var parsedBody loginResult
	err = json.Unmarshal(respBody, &parsedBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body (body: %s): %w", string(respBody), err)
	}

	return &AuthToken{
		Token:     parsedBody.Token,
		ExpiresAt: utils.TimeFromUnix(parsedBody.ExpiresAt),
	}, nil
}
