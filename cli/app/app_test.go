package app

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/tag/config"
	"github.com/FollowTheProcess/tag/git"
)

// setup creates a tempdir, initialises a git repo with a README to do
// search and replace on, adds an initial tag v0.1.0 as well as a tag config
// returns the path to the root of the test repo as well as a teardown
// function that removes the entire directory at the end of the test.
// Usage in a test would be:
//  tmp, teardown := setup(t)
//  defer teardown()
func setup(t *testing.T) (string, func()) {
	t.Helper()
	tmp, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Could not create temp dir: %v", err)
	}
	err = os.WriteFile(filepath.Join(tmp, "README.md"), []byte("Hello, version 0.1.0"), 0755)
	if err != nil {
		t.Fatalf("Could not create README.md: %v", err)
	}

	cfg := []byte(`
	[tag]
	files = [
		{ path = "README.md", search = "version {{.Current}}", replace = "version {{.Next}}" },
	]
	`)

	err = os.WriteFile(filepath.Join(tmp, ".tag.toml"), cfg, 0755)
	if err != nil {
		t.Fatalf("Could not create .tag.toml: %v", err)
	}

	gitConfigEmail := exec.Command("git", "config", "--local", "user.email", "tagtest@gmail.com")
	gitConfigEmail.Dir = tmp

	gitConfigName := exec.Command("git", "config", "--local", "user.name", "Tag Test")
	gitConfigName.Dir = tmp

	init := exec.Command("git", "init", "--initial-branch=main")
	init.Dir = tmp

	add := exec.Command("git", "add", "-A")
	add.Dir = tmp

	commit := exec.Command("git", "commit", "-m", "test commit")
	commit.Dir = tmp

	firstTag := exec.Command("git", "tag", "-a", "v0.1.0", "-m", "test tag")
	firstTag.Dir = tmp

	err = init.Run()
	if err != nil {
		t.Fatalf("Error initialising test git repo: %v", err)
	}

	stdout, err := add.CombinedOutput()
	if err != nil {
		t.Fatalf("Error adding files to test git repo: %s", string(stdout))
	}

	stdout, err = gitConfigEmail.CombinedOutput()
	if err != nil {
		t.Fatalf("git config user.email returned an error: %s", string(stdout))
	}

	stdout, err = gitConfigName.CombinedOutput()
	if err != nil {
		t.Fatalf("git config user.name returned an error: %s", string(stdout))
	}

	stdout, err = commit.CombinedOutput()
	if err != nil {
		t.Fatalf("Error committing to the test git repo: %s", string(stdout))
	}

	stdout, err = firstTag.CombinedOutput()
	if err != nil {
		t.Fatalf("Error issuing the first tag to test git repo: %s", string(stdout))
	}

	tearDown := func() { os.RemoveAll(tmp) }

	return tmp, tearDown
}

// newTestApp creates an app set up for testing.
func newTestApp(out io.Writer) *App {
	app := &App{
		out:     out,
		printer: msg.Default(),
		config: &config.Config{
			Tag: config.Tag{
				Files: []config.File{
					{
						Path:    "README.md",
						Search:  "version {{.Current}}",
						Replace: "version {{.Next}}",
					},
				},
			},
		},
		replace: true,
	}
	return app
}

func TestAppPatch(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	out := &bytes.Buffer{}
	app := newTestApp(out)

	err = app.Patch(true, false, "Message")
	if err != nil {
		t.Fatalf("app.Patch returned an error: %v", err)
	}

	// Check that it's replaced the README contents
	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Could not read from replaced README: %v", err)
	}
	want := "Hello, version 0.1.1"

	if string(readme) != want {
		t.Errorf("README replaced incorrectly: got %q, wanted %q", string(readme), want)
	}

	// Check it's made the appropriate commit
	gitLog := exec.Command("git", "log", "--oneline")
	stdout, err := gitLog.CombinedOutput()
	if err != nil {
		t.Fatalf("Error getting git log: %s", string(stdout))
	}

	if !strings.Contains(string(stdout), "Bump version 0.1.0 -> 0.1.1") {
		t.Errorf("Expected bump version commit not found in git log: %s", string(stdout))
	}

	// Check the working tree is clean
	dirty, err := git.IsDirty()
	if err != nil {
		t.Fatalf("git.IsDirty returned an error: %v", err)
	}
	if dirty {
		t.Error("Working tree was left dirty after replacing files")
	}

	// Check the latest tag is correct
	latest, err := git.LatestTag()
	if err != nil {
		t.Errorf("Could not get latest tag: %v", err)
	}
	if latest != "v0.1.1" {
		t.Errorf("Wrong latest tag: got %s, wanted %s", latest, "v0.1.1")
	}
}

func TestAppMinor(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	out := &bytes.Buffer{}
	app := newTestApp(out)

	err = app.Minor(true, false, "Message")
	if err != nil {
		t.Fatalf("app.Minor returned an error: %v", err)
	}

	// Check that it's replaced the README contents
	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Could not read from replaced README: %v", err)
	}
	want := "Hello, version 0.2.0"

	if string(readme) != want {
		t.Errorf("README replaced incorrectly: got %q, wanted %q", string(readme), want)
	}

	// Check it's made the appropriate commit
	gitLog := exec.Command("git", "log", "--oneline")
	stdout, err := gitLog.CombinedOutput()
	if err != nil {
		t.Fatalf("Error getting git log: %s", string(stdout))
	}

	if !strings.Contains(string(stdout), "Bump version 0.1.0 -> 0.2.0") {
		t.Errorf("Expected bump version commit not found in git log: %s", string(stdout))
	}

	// Check the working tree is clean
	dirty, err := git.IsDirty()
	if err != nil {
		t.Fatalf("git.IsDirty returned an error: %v", err)
	}
	if dirty {
		t.Error("Working tree was left dirty after replacing files")
	}

	// Check the latest tag is correct
	latest, err := git.LatestTag()
	if err != nil {
		t.Errorf("Could not get latest tag: %v", err)
	}
	if latest != "v0.2.0" {
		t.Errorf("Wrong latest tag: got %s, wanted %s", latest, "v0.2.0")
	}
}

func TestAppMajor(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	out := &bytes.Buffer{}
	app := newTestApp(out)

	err = app.Major(true, false, "Message")
	if err != nil {
		t.Fatalf("app.Major returned an error: %v", err)
	}

	// Check that it's replaced the README contents
	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Could not read from replaced README: %v", err)
	}
	want := "Hello, version 1.0.0"

	if string(readme) != want {
		t.Errorf("README replaced incorrectly: got %q, wanted %q", string(readme), want)
	}

	// Check it's made the appropriate commit
	gitLog := exec.Command("git", "log", "--oneline")
	stdout, err := gitLog.CombinedOutput()
	if err != nil {
		t.Fatalf("Error getting git log: %s", string(stdout))
	}

	if !strings.Contains(string(stdout), "Bump version 0.1.0 -> 1.0.0") {
		t.Errorf("Expected bump version commit not found in git log: %s", string(stdout))
	}

	// Check the working tree is clean
	dirty, err := git.IsDirty()
	if err != nil {
		t.Fatalf("git.IsDirty returned an error: %v", err)
	}
	if dirty {
		t.Error("Working tree was left dirty after replacing files")
	}

	// Check the latest tag is correct
	latest, err := git.LatestTag()
	if err != nil {
		t.Errorf("Could not get latest tag: %v", err)
	}
	if latest != "v1.0.0" {
		t.Errorf("Wrong latest tag: got %s, wanted %s", latest, "v1.0.0")
	}
}

func TestAppLatest(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	out := &bytes.Buffer{}
	app := newTestApp(out)

	err = app.Latest()
	if err != nil {
		t.Fatalf("app.Latest returned an error: %v", err)
	}

	if out.String() != "v0.1.0\n" {
		t.Errorf("app.Latest incorrect stdout: got %q, wanted %q", out.String(), "v0.1.0\n")
	}
}

func TestAppList(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	out := &bytes.Buffer{}
	app := newTestApp(out)

	err = app.List()
	if err != nil {
		t.Fatalf("app.List returned an error: %v", err)
	}

	if out.String() != "v0.1.0\n" {
		t.Errorf("app.List incorrect stdout: got %q, wanted %q", out.String(), "v0.1.0\n")
	}
}
