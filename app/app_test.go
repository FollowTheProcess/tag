package app //nolint: testpackage // We need access to some internals

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FollowTheProcess/tag/config"
	"github.com/FollowTheProcess/tag/git"
)

const (
	initialVersion       = "v0.1.0"
	initialReadmeContent = "Hello, version 0.1.0"
	defaultBumpTemplate  = "Bump version {{.Current}} -> {{.Next}}"
	defaultTagTemplate   = "v{{.Next}}"
)

// setup creates a tempdir, initialises a git repo with a README to do
// search and replace on, adds an initial tag v0.1.0 as well as a tag config
// returns the path to the root of the test repo as well as a teardown
// function that removes the entire directory at the end of the test.
// Usage in a test would be:
//
//	tmp, teardown := setup(t)
//	defer teardown()
func setup(t *testing.T) (string, func()) {
	t.Helper()
	tmp, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Could not create temp dir: %v", err)
	}
	err = os.WriteFile(filepath.Join(tmp, "README.md"), []byte(initialReadmeContent), 0o755)
	if err != nil {
		t.Fatalf("Could not create README.md: %v", err)
	}

	cfg := []byte(`
	version = '0.1.0'

	[git]
	default-branch = 'main'
	message-template = 'Bump version {{.Current}} -> {{.Next}}'
	tag-template = 'v{{.Next}}'
	
	[hooks]
	pre-replace = "echo 'I run before doing anything'"
	pre-commit = "echo 'I run after replacing but before committing changes'"
	pre-tag = "echo 'I run after committing changes but before tagging'"
	pre-push = "echo 'I run after tagging, but before pushing'"

	[[file]]
	path = 'README.md'
	search = 'Hello, version {{.Current}}'	
	`)

	err = os.WriteFile(filepath.Join(tmp, ".tag.toml"), cfg, 0o755)
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

	firstTag := exec.Command("git", "tag", "-a", initialVersion, "-m", "test tag")
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
func newTestApp(out io.Writer) App {
	app := App{
		Stdout: out,
		Cfg: config.Config{
			Version: "0.1.0",
			Files: []config.File{
				{
					Path:   "README.md",
					Search: "version = {{.Current}}",
				},
			},
		},
		replaceMode: true,
	}
	return app
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

	if out.String() != fmt.Sprintln(initialVersion) {
		t.Errorf("app.Latest incorrect stdout: got %q, wanted %q", out.String(), fmt.Sprintln(initialVersion))
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

	err = app.List(10)
	if err != nil {
		t.Fatalf("app.List returned an error: %v", err)
	}

	if out.String() != fmt.Sprintln(initialVersion) {
		t.Errorf("app.List incorrect stdout: got %q, wanted %q", out.String(), fmt.Sprintln(initialVersion))
	}
}

func TestAppMajor(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	appOut := &bytes.Buffer{}
	appErr := &bytes.Buffer{}
	app, err := New(tmp, appOut, appErr)
	if err != nil {
		t.Fatalf("app.New returned an error: %v", err)
	}

	err = app.Major(false, true, false)
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

	// Check it's replaced the config version but not the message templates etc.
	cfg, err := config.Load(filepath.Join(tmp, ".tag.toml"))
	if err != nil {
		t.Fatalf("Could not read replaced config file: %v", err)
	}

	if cfg.Version != "1.0.0" {
		t.Errorf("Wrong version in replaced config file. Got %s, wanted %s", cfg.Version, "1.0.0")
	}

	if cfg.Git.MessageTemplate != defaultBumpTemplate {
		t.Errorf("Wrong message template in replaced config file. Got %s, wanted %s", cfg.Git.MessageTemplate, defaultBumpTemplate)
	}

	if cfg.Git.TagTemplate != defaultTagTemplate {
		t.Errorf("Wrong tag template in replaced config file. Got %s, wanted %s", cfg.Git.TagTemplate, defaultTagTemplate)
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

func TestAppMajorDryRun(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	appOut := &bytes.Buffer{}
	appErr := &bytes.Buffer{}
	app, err := New(tmp, appOut, appErr)
	if err != nil {
		t.Fatalf("app.New returned an error: %v", err)
	}

	err = app.Major(false, true, true)
	if err != nil {
		t.Fatalf("app.Major returned an error: %v", err)
	}

	// Check that it's not replaced the README contents
	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Could not read from replaced README: %v", err)
	}

	if string(readme) != initialReadmeContent {
		t.Errorf("README replaced despite dry-run: got %q, wanted %q", string(readme), initialReadmeContent)
	}

	// Check it's not made a commit
	gitLog := exec.Command("git", "log", "--oneline")
	stdout, err := gitLog.CombinedOutput()
	if err != nil {
		t.Fatalf("Error getting git log: %s", string(stdout))
	}

	if !strings.Contains(string(stdout), "test commit") {
		t.Errorf("Made a commit despite dry-run: %s", string(stdout))
	}

	// Check the working tree is clean
	dirty, err := git.IsDirty()
	if err != nil {
		t.Fatalf("git.IsDirty returned an error: %v", err)
	}
	if dirty {
		t.Error("Working tree is dirty after dry-run")
	}

	// Check the latest tag is correct
	latest, err := git.LatestTag()
	if err != nil {
		t.Errorf("Could not get latest tag: %v", err)
	}
	if latest != initialVersion {
		t.Errorf("Wrong latest tag: got %s, wanted %s", latest, initialVersion)
	}
}

func TestAppMinor(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	appOut := &bytes.Buffer{}
	appErr := &bytes.Buffer{}
	app, err := New(tmp, appOut, appErr)
	if err != nil {
		t.Fatalf("app.New returned an error: %v", err)
	}

	err = app.Minor(false, true, false)
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

	// Check it's replaced the config version but not the message templates etc.
	cfg, err := config.Load(filepath.Join(tmp, ".tag.toml"))
	if err != nil {
		t.Fatalf("Could not read replaced config file: %v", err)
	}

	if cfg.Version != "0.2.0" {
		t.Errorf("Wrong version in replaced config file. Got %s, wanted %s", cfg.Version, "0.2.0")
	}

	if cfg.Git.MessageTemplate != defaultBumpTemplate {
		t.Errorf("Wrong message template in replaced config file. Got %s, wanted %s", cfg.Git.MessageTemplate, defaultBumpTemplate)
	}

	if cfg.Git.TagTemplate != defaultTagTemplate {
		t.Errorf("Wrong tag template in replaced config file. Got %s, wanted %s", cfg.Git.TagTemplate, defaultTagTemplate)
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

func TestAppMinorDryRun(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	appOut := &bytes.Buffer{}
	appErr := &bytes.Buffer{}
	app, err := New(tmp, appOut, appErr)
	if err != nil {
		t.Fatalf("app.New returned an error: %v", err)
	}

	err = app.Minor(false, true, true)
	if err != nil {
		t.Fatalf("app.Minor returned an error: %v", err)
	}

	// Check that it's not replaced the README contents
	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Could not read from replaced README: %v", err)
	}

	if string(readme) != initialReadmeContent {
		t.Errorf("README replaced despite dry-run: got %q, wanted %q", string(readme), initialReadmeContent)
	}

	// Check it's not made a commit
	gitLog := exec.Command("git", "log", "--oneline")
	stdout, err := gitLog.CombinedOutput()
	if err != nil {
		t.Fatalf("Error getting git log: %s", string(stdout))
	}

	if !strings.Contains(string(stdout), "test commit") {
		t.Errorf("Made a commit despite dry-run: %s", string(stdout))
	}

	// Check the working tree is clean
	dirty, err := git.IsDirty()
	if err != nil {
		t.Fatalf("git.IsDirty returned an error: %v", err)
	}
	if dirty {
		t.Error("Working tree is dirty after dry-run")
	}

	// Check the latest tag is correct
	latest, err := git.LatestTag()
	if err != nil {
		t.Errorf("Could not get latest tag: %v", err)
	}
	if latest != initialVersion {
		t.Errorf("Wrong latest tag: got %s, wanted %s", latest, initialVersion)
	}
}

func TestAppPatch(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	appOut := &bytes.Buffer{}
	appErr := &bytes.Buffer{}
	app, err := New(tmp, appOut, appErr)
	if err != nil {
		t.Fatalf("app.New returned an error: %v", err)
	}

	err = app.Patch(false, true, false)
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

	// Check it's replaced the config version but not the message templates etc.
	cfg, err := config.Load(filepath.Join(tmp, ".tag.toml"))
	if err != nil {
		t.Fatalf("Could not read replaced config file: %v", err)
	}

	if cfg.Version != "0.1.1" {
		t.Errorf("Wrong version in replaced config file. Got %s, wanted %s", cfg.Version, "0.1.1")
	}

	if cfg.Git.MessageTemplate != defaultBumpTemplate {
		t.Errorf("Wrong message template in replaced config file. Got %s, wanted %s", cfg.Git.MessageTemplate, defaultBumpTemplate)
	}

	if cfg.Git.TagTemplate != defaultTagTemplate {
		t.Errorf("Wrong tag template in replaced config file. Got %s, wanted %s", cfg.Git.TagTemplate, defaultTagTemplate)
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

func TestAppPatchDryRun(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	appOut := &bytes.Buffer{}
	appErr := &bytes.Buffer{}
	app, err := New(tmp, appOut, appErr)
	if err != nil {
		t.Fatalf("app.New returned an error: %v", err)
	}

	err = app.Patch(false, true, true)
	if err != nil {
		t.Fatalf("app.Patch returned an error: %v", err)
	}

	// Check that it's not replaced the README contents
	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Could not read from replaced README: %v", err)
	}

	if string(readme) != initialReadmeContent {
		t.Errorf("README replaced despite dry-run: got %q, wanted %q", string(readme), initialReadmeContent)
	}

	// Check it's not made a commit
	gitLog := exec.Command("git", "log", "--oneline")
	stdout, err := gitLog.CombinedOutput()
	if err != nil {
		t.Fatalf("Error getting git log: %s", string(stdout))
	}

	if !strings.Contains(string(stdout), "test commit") {
		t.Errorf("Made a commit despite dry-run: %s", string(stdout))
	}

	// Check the working tree is clean
	dirty, err := git.IsDirty()
	if err != nil {
		t.Fatalf("git.IsDirty returned an error: %v", err)
	}
	if dirty {
		t.Error("Working tree is dirty after dry-run")
	}

	// Check the latest tag is correct
	latest, err := git.LatestTag()
	if err != nil {
		t.Errorf("Could not get latest tag: %v", err)
	}
	if latest != initialVersion {
		t.Errorf("Wrong latest tag: got %s, wanted %s", latest, initialVersion)
	}
}
