// Package github fetches a user's public profile and repositories from the
// GitHub REST API.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// DefaultBaseURL is the root of the public GitHub REST API.
const DefaultBaseURL = "https://api.github.com"

const perPage = 100

// User is the subset of the GitHub user object the site uses.
type User struct {
	Login           string `json:"login"`
	Name            string `json:"name"`
	AvatarURL       string `json:"avatar_url"`
	HTMLURL         string `json:"html_url"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	Email           string `json:"email"`
	Bio             string `json:"bio"`
	TwitterUsername string `json:"twitter_username"`
	PublicRepos     int    `json:"public_repos"`
}

// Repo is the subset of the GitHub repository object the site uses.
type Repo struct {
	Name            string    `json:"name"`
	HTMLURL         string    `json:"html_url"`
	Description     string    `json:"description"`
	Homepage        string    `json:"homepage"`
	Language        string    `json:"language"`
	Topics          []string  `json:"topics"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	Fork            bool      `json:"fork"`
	Archived        bool      `json:"archived"`
	PushedAt        time.Time `json:"pushed_at"`
}

// Client talks to the GitHub REST API.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient returns a client for the public GitHub API. The token is
// optional; without one GitHub allows 60 requests per hour per IP.
func NewClient(token string) *Client {
	return &Client{
		BaseURL:    DefaultBaseURL,
		Token:      token,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// User fetches the public profile of login.
func (c *Client) User(ctx context.Context, login string) (User, error) {
	var user User
	err := c.get(ctx, "/users/"+url.PathEscape(login), nil, &user)
	return user, err
}

// Repos fetches every public repository owned by login, following pagination.
func (c *Client) Repos(ctx context.Context, login string) ([]Repo, error) {
	var repos []Repo
	for page := 1; ; page++ {
		query := url.Values{
			"type":     {"owner"},
			"sort":     {"pushed"},
			"per_page": {strconv.Itoa(perPage)},
			"page":     {strconv.Itoa(page)},
		}
		var batch []Repo
		if err := c.get(ctx, "/users/"+url.PathEscape(login)+"/repos", query, &batch); err != nil {
			return nil, err
		}
		repos = append(repos, batch...)
		if len(batch) < perPage {
			return repos, nil
		}
	}
}

func (c *Client) get(ctx context.Context, path string, query url.Values, v any) error {
	endpoint := c.BaseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr struct {
			Message string `json:"message"`
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Message != "" {
			return fmt.Errorf("GET %s: %s: %s", path, resp.Status, apiErr.Message)
		}
		return fmt.Errorf("GET %s: %s", path, resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("GET %s: decoding response: %w", path, err)
	}
	return nil
}
