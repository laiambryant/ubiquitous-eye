// Package site renders the portfolio page and writes it, together with its
// static assets, to an output directory.
package site

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"time"

	"github.com/laiambryant/ubiquitous-eye/internal/github"
)

//go:embed templates/index.html assets
var files embed.FS

var (
	assetVersion = hashAssets()
	tmpl         = template.Must(template.New("index.html").Funcs(template.FuncMap{
		"asset":     assetURL,
		"date":      func(t time.Time) string { return t.Format("Jan 2, 2006") },
		"isoDate":   func(t time.Time) string { return t.Format("2006-01-02") },
		"langColor": languageColor,
	}).ParseFS(files, "templates/index.html"))
)

// Page is the data rendered into the site.
type Page struct {
	User  github.User
	Repos []github.Repo
	// Role is shown above the name, e.g. "Software Developer". Optional.
	Role string
	// Email is used for the contact links. Optional.
	Email string
	Year  int
}

// Language is one row of the language breakdown.
type Language struct {
	Name  string
	Color template.CSS
	Repos int
}

// DisplayName is the user's name, falling back to their login.
func (p Page) DisplayName() string {
	return cmp.Or(p.User.Name, p.User.Login)
}

// Description is used for the page's meta description.
func (p Page) Description() string {
	return cmp.Or(p.User.Bio, "Open source projects by "+p.DisplayName()+".")
}

// Languages counts repositories by primary language, most common first.
func (p Page) Languages() []Language {
	counts := map[string]int{}
	for _, r := range p.Repos {
		if r.Language != "" {
			counts[r.Language]++
		}
	}
	langs := make([]Language, 0, len(counts))
	for name, n := range counts {
		langs = append(langs, Language{Name: name, Color: languageColor(name), Repos: n})
	}
	slices.SortFunc(langs, func(a, b Language) int {
		return cmp.Or(cmp.Compare(b.Repos, a.Repos), cmp.Compare(a.Name, b.Name))
	})
	return langs
}

// Render writes the HTML page to w. Repositories are listed most recently
// pushed first.
func Render(w io.Writer, p Page) error {
	p.Repos = slices.Clone(p.Repos)
	slices.SortStableFunc(p.Repos, func(a, b github.Repo) int {
		return b.PushedAt.Compare(a.PushedAt)
	})
	return tmpl.Execute(w, p)
}

// Generate writes index.html and its assets into dir, creating it if needed.
func Generate(dir string, p Page) error {
	var buf bytes.Buffer
	if err := Render(&buf, p); err != nil {
		return err
	}
	if err := writeAssets(dir); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "index.html"), buf.Bytes(), 0o644)
}

func writeAssets(dir string) error {
	return fs.WalkDir(files, "assets", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		target := filepath.Join(dir, filepath.FromSlash(name))
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := files.ReadFile(name)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

// hashAssets fingerprints the embedded assets so browsers refetch them only
// when they change.
func hashAssets() string {
	h := sha256.New()
	err := fs.WalkDir(files, "assets", func(name string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := files.ReadFile(name)
		h.Write(data)
		return err
	})
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil))[:10]
}

func assetURL(name string) string {
	return path.Join("assets", name) + "?v=" + assetVersion
}

// Colors match GitHub Linguist so the language dots look familiar.
var languageColors = map[string]template.CSS{
	"C":          "#555555",
	"C#":         "#178600",
	"C++":        "#f34b7d",
	"CSS":        "#663399",
	"Dart":       "#00b4ab",
	"Dockerfile": "#384d54",
	"Elixir":     "#6e4a7e",
	"Go":         "#00add8",
	"HTML":       "#e34c26",
	"Haskell":    "#5e5086",
	"Java":       "#b07219",
	"JavaScript": "#f1e05a",
	"Kotlin":     "#a97bff",
	"Lua":        "#000080",
	"Nix":        "#7e7eff",
	"PHP":        "#4f5d95",
	"Python":     "#3572a5",
	"Ruby":       "#701516",
	"Rust":       "#dea584",
	"Shell":      "#89e051",
	"Swift":      "#f05138",
	"TypeScript": "#3178c6",
	"Vim Script": "#199f4b",
	"Zig":        "#ec915c",
}

func languageColor(name string) template.CSS {
	if c, ok := languageColors[name]; ok {
		return c
	}
	return "#8b949e"
}
