package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const defaultURL = "https://api.github.com"

type Client struct {
	HTTP    *http.Client
	BaseURL string
	Token   string
}

func NewClient() *Client {
	return &Client{
		HTTP:    &http.Client{Timeout: 10 * time.Second},
		BaseURL: defaultURL,
		Token:   os.Getenv("GITHUB_TOKEN"),
	}
}

func (c *Client) GetUser(ctx context.Context, username string) (User, error) {
	var u User
	err := c.getJSON(ctx, fmt.Sprintf("%s/users/%s", c.BaseURL, url.PathEscape(username)), &u)
	return u, err
}

func (c *Client) GetRepos(ctx context.Context, username string) ([]Repo, error) {
	var all []Repo
	for page := 1; ; page++ {
		endpoint := fmt.Sprintf("%s/users/%s/repos?per_page=100&page=%d",
			c.BaseURL, url.PathEscape(username), page)

		var batch []Repo
		if err := c.getJSON(ctx, endpoint, &batch); err != nil {
			return nil, err
		}
		all = append(all, batch...)
		if len(batch) < 100 {
			return all, nil
		}
	}
}

func (c *Client) GetLanguages(ctx context.Context, username, repo string) (Languages, error) {
	var l Languages
	endpoint := fmt.Sprintf("%s/repos/%s/%s/languages", c.BaseURL, url.PathEscape(username), url.PathEscape(repo))
	err := c.getJSON(ctx, endpoint, &l)
	return l, err
}

func (c *Client) getJSON(ctx context.Context, endpoint string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", endpoint, resp.Status)
	}
	return json.Unmarshal(body, v)
}
