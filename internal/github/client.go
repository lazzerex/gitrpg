package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const apiURL = "https://api.github.com/graphql"

var ErrUnauthorized = errors.New("github: access token rejected")

func statusError(status int) error {
	if status == http.StatusUnauthorized {
		return fmt.Errorf("github graphql: HTTP 401: %w", ErrUnauthorized)
	}
	return fmt.Errorf("github graphql: HTTP %d", status)
}

type Client struct {
	http       *http.Client
	token      string
	logger     *slog.Logger
	PointsUsed int
}

func newClient(token string, logger *slog.Logger) *Client {
	return &Client{http: &http.Client{Timeout: 30 * time.Second}, token: token, logger: logger}
}

type gqlReq struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type gqlResp struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *Client) query(ctx context.Context, q string, vars map[string]any, out any) error {
	payload, err := json.Marshal(gqlReq{Query: q, Variables: vars})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return statusError(resp.StatusCode)
	}

	var gr gqlResp
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return err
	}
	if len(gr.Errors) > 0 {
		return fmt.Errorf("github graphql: %s", gr.Errors[0].Message)
	}

	return json.Unmarshal(gr.Data, out)
}
