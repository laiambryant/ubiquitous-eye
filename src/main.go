// Command ubiquitous-eye builds a static portfolio site from a GitHub user's
// public profile and repositories.
package main

import (
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/laiambryant/ubiquitous-eye/internal/github"
	"github.com/laiambryant/ubiquitous-eye/internal/site"
)

type options struct {
	user     string
	email    string
	role     string
	out      string
	logLevel slog.Level
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := github.NewClient(os.Getenv("GITHUB_TOKEN"))
	err := run(ctx, os.Args[1:], os.Stderr, client)
	switch {
	case err == nil:
	case errors.Is(err, flag.ErrHelp):
		os.Exit(0)
	case errors.As(err, new(usageError)):
		os.Exit(2)
	default:
		fmt.Fprintln(os.Stderr, "ubiquitous-eye:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stderr io.Writer, client *github.Client) error {
	opts, err := parseFlags(args, stderr)
	if err != nil {
		return err
	}
	logger := slog.New(slog.NewTextHandler(stderr, &slog.HandlerOptions{Level: opts.logLevel}))

	logger.Debug("fetching profile", "user", opts.user)
	user, err := client.User(ctx, opts.user)
	if err != nil {
		return fmt.Errorf("fetching profile for %q: %w", opts.user, err)
	}
	logger.Debug("fetching repositories", "user", opts.user)
	repos, err := client.Repos(ctx, opts.user)
	if err != nil {
		return fmt.Errorf("fetching repositories for %q: %w", opts.user, err)
	}

	page := site.Page{
		User:  user,
		Repos: repos,
		Role:  opts.role,
		Email: cmp.Or(opts.email, user.Email),
		Year:  time.Now().Year(),
	}
	if err := site.Generate(opts.out, page); err != nil {
		return fmt.Errorf("writing site to %s: %w", opts.out, err)
	}
	logger.Info("site generated", "user", opts.user, "repos", len(repos), "out", opts.out)
	return nil
}

type usageError struct{ msg string }

func (e usageError) Error() string { return e.msg }

func parseFlags(args []string, output io.Writer) (options, error) {
	opts := options{logLevel: slog.LevelInfo}
	fs := flag.NewFlagSet("ubiquitous-eye", flag.ContinueOnError)
	fs.SetOutput(output)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: ubiquitous-eye -user <github-username> [flags]\n\n")
		fmt.Fprintf(fs.Output(), "Builds a static portfolio site from a GitHub user's public profile and repositories.\n")
		fmt.Fprintf(fs.Output(), "Set GITHUB_TOKEN to raise the GitHub API rate limit.\n\nFlags:\n")
		fs.PrintDefaults()
	}
	fs.StringVar(&opts.user, "user", "", "GitHub username to build the site for (required)")
	fs.StringVar(&opts.out, "out", "docs", "directory to write the site to")
	fs.StringVar(&opts.email, "email", "", "contact email (defaults to the public email on the GitHub profile)")
	fs.StringVar(&opts.role, "role", "", `role shown above the name, e.g. "Software Developer"`)
	fs.TextVar(&opts.logLevel, "log-level", slog.LevelInfo, "log level: debug, info, warn or error")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return opts, err
		}
		return opts, usageError{err.Error()}
	}
	if fs.NArg() > 0 {
		return opts, usage(fs, "unexpected arguments: %v", fs.Args())
	}
	if opts.user == "" {
		return opts, usage(fs, "-user is required")
	}
	if opts.out == "" {
		return opts, usage(fs, "-out must not be empty")
	}
	return opts, nil
}

func usage(fs *flag.FlagSet, format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintln(fs.Output(), msg)
	fs.Usage()
	return usageError{msg}
}
