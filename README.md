# Ubiquitous Eye 👁️

[![Go Reference](https://pkg.go.dev/badge/github.com/laiambryant/ubiquitous-eye.svg)](https://pkg.go.dev/github.com/laiambryant/ubiquitous-eye)
[![Go Report Card](https://goreportcard.com/badge/github.com/laiambryant/ubiquitous-eye)](https://goreportcard.com/report/github.com/laiambryant/ubiquitous-eye)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/laiambryant/ubiquitous-eye/workflows/CI/badge.svg)](https://github.com/laiambryant/ubiquitous-eye/actions)

A simple Go application that generates a static portfolio website from GitHub user data. Fetches your GitHub profile and repositories to create a responsive portfolio page with light/dark mode theming.

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher

### Installation & Usage

```bash
# Clone the repository
git clone https://github.com/laiambryant/ubiquitous-eye.git
cd ubiquitous-eye/src

# Install dependencies
go mod tidy

# Generate portfolio site
go run main.go
```

This creates a static `docs/index.html` file ready for deployment.

## What It Does

The application:

1. Fetches your GitHub user profile from the GitHub API
2. Retrieves your public repositories
3. Renders this data into a responsive HTML portfolio template
4. Outputs a static `docs/index.html` file

The generated portfolio includes:

- Your GitHub profile information (name, bio, avatar, etc.)
- A grid of your public repositories
- Light/dark mode toggle
- Responsive design for mobile and desktop

## Configuration

To customize for a different GitHub user, update the username in `src/packages/utils/constants.go`:

```go
const GITHUB_USER_LBRYANT = "https://api.github.com/users/YOUR_USERNAME"
```

## Project Structure

```text
ubiquitous-eye/
├── src/
│   ├── main.go                           # Entry point
│   ├── go.mod                            # Go module
│   ├── packages/
│   │   ├── model/response/               # GitHub API response structs
│   │   ├── services/                     # Core logic
│   │   │   ├── services.go               # Data fetching & rendering
│   │   │   └── resources/index.html      # Portfolio template
│   │   └── utils/constants.go            # Configuration
└── docs/                                 # Generated output (index.html)
```

## Deployment

The generated `docs/index.html` can be deployed to:

- GitHub Pages (set source to `/docs` folder)
- Any static hosting service
- Local web server

---

⭐ **Star this repository if you find it helpful!**