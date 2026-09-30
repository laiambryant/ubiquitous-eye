package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := NewClient("")
	c.BaseURL = srv.URL
	return c
}

func TestUser(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/octocat" {
			http.NotFound(w, r)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer secret")
		}
		fmt.Fprint(w, `{"login":"octocat","name":"The Octocat","location":"San Francisco","public_repos":8}`)
	}))
	c.Token = "secret"

	user, err := c.User(context.Background(), "octocat")
	if err != nil {
		t.Fatal(err)
	}
	if user.Name != "The Octocat" || user.Location != "San Francisco" || user.PublicRepos != 8 {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestUserOmitsAuthorizationWithoutToken(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("Authorization header = %q, want none", got)
		}
		fmt.Fprint(w, `{"login":"octocat"}`)
	}))

	if _, err := c.User(context.Background(), "octocat"); err != nil {
		t.Fatal(err)
	}
}

func TestReposFollowsPagination(t *testing.T) {
	const total = perPage + 3
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("sort") != "pushed" || q.Get("type") != "owner" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		page, _ := strconv.Atoi(q.Get("page"))
		start := (page - 1) * perPage
		var repos []Repo
		for i := start; i < total && i < start+perPage; i++ {
			repos = append(repos, Repo{Name: "repo-" + strconv.Itoa(i)})
		}
		if repos == nil {
			repos = []Repo{}
		}
		json.NewEncoder(w).Encode(repos)
	}))

	repos, err := c.Repos(context.Background(), "octocat")
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != total {
		t.Fatalf("got %d repos, want %d", len(repos), total)
	}
	if repos[total-1].Name != "repo-"+strconv.Itoa(total-1) {
		t.Errorf("last repo = %q", repos[total-1].Name)
	}
}

func TestErrorIncludesAPIMessage(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"message":"API rate limit exceeded"}`)
	}))

	_, err := c.User(context.Background(), "octocat")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "403") || !strings.Contains(err.Error(), "API rate limit exceeded") {
		t.Errorf("error %q should mention the status and the API message", err)
	}
}

func TestUserNotFound(t *testing.T) {
	c := newTestClient(t, http.NotFoundHandler())

	if _, err := c.User(context.Background(), "nobody"); err == nil {
		t.Fatal("expected an error for a missing user")
	}
}
