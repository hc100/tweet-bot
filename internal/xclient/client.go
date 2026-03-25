package xclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hc100/tweet-bot/internal/config"
)

const createTweetURL = "https://api.x.com/2/tweets"

type Client struct {
	httpClient  *http.Client
	credentials config.Credentials
}

func NewClient(credentials config.Credentials, timeout time.Duration) *Client {
	return &Client{
		httpClient:  &http.Client{Timeout: timeout},
		credentials: credentials,
	}
}

func (c *Client) Post(ctx context.Context, text string) error {
	payload, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createTweetURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", buildAuthorizationHeader(
		http.MethodPost,
		createTweetURL,
		nil,
		c.credentials,
	))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("x api returned status=%d body=%s", resp.StatusCode, string(body))
	}

	return nil
}
