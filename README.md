# ubiquitous-eye

[![CI](https://github.com/laiambryant/ubiquitous-eye/actions/workflows/ci.yml/badge.svg)](https://github.com/laiambryant/ubiquitous-eye/actions/workflows/ci.yml)
[![Deploy site](https://github.com/laiambryant/ubiquitous-eye/actions/workflows/build-deploy.yml/badge.svg)](https://github.com/laiambryant/ubiquitous-eye/actions/workflows/build-deploy.yml)

Generates a static portfolio page from a GitHub user's public profile and repositories.

Live site: <https://laiambryant.github.io/ubiquitous-eye/>

## Usage

Requires Go 1.26 or newer.

```sh
cd src
go run . -user laiambryant -email liambryant@outlook.it -role "Software Developer" -out ../docs
```

This writes `index.html` and `assets/` into the output directory.

| Flag         | Default | Description                                                         |
| ------------ | ------- | ------------------------------------------------------------------- |
| `-user`      |         | GitHub username to build the site for. Required.                    |
| `-out`       | `docs`  | Directory to write the site to.                                     |
| `-email`     |         | Contact email. Falls back to the public email on the GitHub profile. |
| `-role`      |         | Line shown above the name, e.g. `Software Developer`.               |
| `-log-level` | `info`  | `debug`, `info`, `warn` or `error`.                                 |

Unauthenticated GitHub API requests are limited to 60 per hour. Set `GITHUB_TOKEN` to use a token instead:

```sh
GITHUB_TOKEN=$(gh auth token) go run . -user laiambryant
```

Repositories are listed most recently pushed first. Name, bio, location, company and avatar come from the GitHub profile, so edit them there.

## Development

```sh
cd src
go test ./...
```

The tests use a local HTTP server instead of the GitHub API, so they run offline.

```text
src/
├── main.go                 CLI flags and wiring
└── internal/
    ├── github/             GitHub REST API client
    └── site/               HTML rendering
        ├── templates/      page template
        └── assets/         CSS and JS copied next to index.html
```

## Deployment

GitHub Pages serves `/docs` from the `deploy` branch. The [Deploy site](.github/workflows/build-deploy.yml) workflow runs on every push to `main`, once a day, and on manual dispatch. It runs the tests, generates the site and commits it to `deploy`. If nothing changed, it skips the commit.

To build a site for another account, change the flags in that workflow.
