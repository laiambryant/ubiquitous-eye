package site

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/laiambryant/ubiquitous-eye/internal/github"
)

func testPage() Page {
	day := func(d int) time.Time { return time.Date(2026, time.September, d, 12, 0, 0, 0, time.UTC) }
	return Page{
		User: github.User{
			Login:     "octocat",
			Name:      "The Octocat",
			AvatarURL: "https://avatars.example.com/octocat.png",
			HTMLURL:   "https://github.com/octocat",
			Location:  "San Francisco",
			Bio:       "Builds things.",
		},
		Repos: []github.Repo{
			{Name: "older", Language: "Go", PushedAt: day(1)},
			{Name: "newest", Language: "C", Description: "<script>alert(1)</script>", StargazersCount: 3, PushedAt: day(14)},
			{Name: "middle", Language: "Go", Topics: []string{"cli"}, Homepage: "https://example.com", PushedAt: day(7)},
		},
		Role:  "Software Developer",
		Email: "octo@example.com",
		Year:  2026,
	}
}

func render(t *testing.T, p Page) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Render(&buf, p); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestRender(t *testing.T) {
	out := render(t, testPage())

	for _, want := range []string{
		"<title>The Octocat · Software Developer</title>",
		`href="mailto:octo@example.com"`,
		`<time datetime="2026-09-14">Sep 14, 2026</time>`,
		`<a class="repo-site" href="https://example.com">Website</a>`,
		"3 public repositories",
		"&lt;script&gt;",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output is missing %q", want)
		}
	}
	for _, unwanted := range []string{"<script>alert(1)</script>", "ZgotmplZ"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("output contains %q", unwanted)
		}
	}
}

func TestRenderOrdersReposByLastPush(t *testing.T) {
	out := render(t, testPage())

	newest := strings.Index(out, ">newest<")
	middle := strings.Index(out, ">middle<")
	older := strings.Index(out, ">older<")
	if newest < 0 || middle < 0 || older < 0 {
		t.Fatal("output is missing a repository")
	}
	if !(newest < middle && middle < older) {
		t.Errorf("repositories are not ordered by last push: newest=%d middle=%d older=%d", newest, middle, older)
	}
}

func TestRenderOmitsOptionalFields(t *testing.T) {
	p := testPage()
	p.User.Name = ""
	p.Role = ""
	p.Email = ""
	out := render(t, p)

	if !strings.Contains(out, "<title>octocat</title>") {
		t.Error("title should fall back to the login when the name is empty")
	}
	if strings.Contains(out, "mailto:") {
		t.Error("output should not contain an email link when no email is set")
	}
	if strings.Contains(out, `class="role"`) {
		t.Error("output should not contain a role when none is set")
	}
}

func TestLanguages(t *testing.T) {
	p := testPage()
	p.Repos = append(p.Repos, github.Repo{Name: "docs"})

	got := p.Languages()
	if len(got) != 2 {
		t.Fatalf("got %d languages, want 2: %+v", len(got), got)
	}
	if got[0].Name != "Go" || got[0].Repos != 2 || got[0].Color != "#00add8" {
		t.Errorf("first language = %+v, want Go with 2 repos", got[0])
	}
	if got[1].Name != "C" || got[1].Repos != 1 {
		t.Errorf("second language = %+v, want C with 1 repo", got[1])
	}
}

func TestGenerate(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "docs")

	if err := Generate(dir, testPage()); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"index.html", "assets/site.css", "assets/site.js"} {
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil {
			t.Errorf("expected %s to be written: %v", name, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("%s is empty", name)
		}
	}
}
