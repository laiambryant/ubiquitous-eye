package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/laiambryant/ubiquitous-eye/internal/github"
)

func TestParseFlags(t *testing.T) {
	opts, err := parseFlags([]string{"-user", "octocat", "-email", "octo@example.com", "-role", "Developer", "-out", "site", "-log-level", "debug"}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if opts.user != "octocat" || opts.email != "octo@example.com" || opts.role != "Developer" || opts.out != "site" || opts.logLevel.String() != "DEBUG" {
		t.Errorf("unexpected options: %+v", opts)
	}
}

func TestParseFlagsDefaults(t *testing.T) {
	opts, err := parseFlags([]string{"-user", "octocat"}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if opts.out != "docs" || opts.logLevel.String() != "INFO" {
		t.Errorf("unexpected defaults: %+v", opts)
	}
}

func TestParseFlagsErrors(t *testing.T) {
	tests := map[string][]string{
		"missing user":     {},
		"empty out":        {"-user", "octocat", "-out", ""},
		"extra arguments":  {"-user", "octocat", "extra"},
		"unknown flag":     {"-user", "octocat", "-nope"},
		"invalid loglevel": {"-user", "octocat", "-log-level", "loud"},
	}
	for name, args := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := parseFlags(args, io.Discard)
			if !errors.As(err, new(usageError)) {
				t.Errorf("got %v, want a usage error", err)
			}
		})
	}
}

func TestParseFlagsHelp(t *testing.T) {
	var out bytes.Buffer
	_, err := parseFlags([]string{"-h"}, &out)
	if !errors.Is(err, flag.ErrHelp) {
		t.Errorf("got %v, want flag.ErrHelp", err)
	}
	if !strings.Contains(out.String(), "-user") {
		t.Errorf("help output should document -user:\n%s", out.String())
	}
}

func TestRun(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users/octocat":
			fmt.Fprint(w, `{"login":"octocat","name":"The Octocat","html_url":"https://github.com/octocat","email":"public@example.com"}`)
		case "/users/octocat/repos":
			fmt.Fprint(w, `[{"name":"hello-world","language":"Go","pushed_at":"2026-09-01T10:00:00Z"}]`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	client := github.NewClient("")
	client.BaseURL = srv.URL
	out := filepath.Join(t.TempDir(), "docs")

	if err := run(context.Background(), []string{"-user", "octocat", "-out", out}, io.Discard, client); err != nil {
		t.Fatal(err)
	}

	html, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"The Octocat", "hello-world", "mailto:public@example.com"} {
		if !bytes.Contains(html, []byte(want)) {
			t.Errorf("index.html is missing %q", want)
		}
	}
}

func TestRunReportsAPIErrors(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	client := github.NewClient("")
	client.BaseURL = srv.URL

	err := run(context.Background(), []string{"-user", "nobody", "-out", t.TempDir()}, io.Discard, client)
	if err == nil || !strings.Contains(err.Error(), "nobody") {
		t.Errorf("got %v, want an error naming the user", err)
	}
}
