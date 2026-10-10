// Package lumaservice is growth-portal's client for luma-service, the only way this app gets Luma data
// (docs/pcr/0003).
package lumaservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	secret  string
	http    *http.Client
}

func New(baseURL, secret string) *Client {
	return &Client{
		baseURL: baseURL,
		secret:  secret,
		// Generous: luma-service pages through Luma for large tiers.
		http: &http.Client{Timeout: 60 * time.Second},
	}
}

// Member is a Community tier member as this app needs it. Names and emails are deliberately not
// decoded: the app only works with aggregates (docs/pcr/0006).
type Member struct {
	CreatedAt time.Time `json:"created_at"`
	Answers   []Answer  `json:"answers"`
}

type Answer struct {
	QuestionID string   `json:"question_id"`
	Label      string   `json:"label"`
	Type       string   `json:"type"`
	Values     []string `json:"values"`
}

func (c *Client) CommunityMembers(ctx context.Context) ([]Member, error) {
	var body struct {
		Members []Member `json:"members"`
	}
	if err := c.get(ctx, "/v1/community/members", &body); err != nil {
		return nil, err
	}
	return body.Members, nil
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Service-Secret", c.secret)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("luma-service: %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("luma-service: %s: status %d: %s", path, resp.StatusCode, body)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("luma-service: %s: decoding response: %w", path, err)
	}
	return nil
}
