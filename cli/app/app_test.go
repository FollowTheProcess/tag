package app

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/tag/config"
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
		{ path = "README.md", search = "version = {{.Current}}", replace = "version = {{.Next}}" },
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

	init := exec.Command("git", "init")
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

func TestAppPatch(t *testing.T) {
	tmp, teardown := setup(t)
	defer teardown()

	err := os.Chdir(tmp)
	if err != nil {
		t.Fatalf("Could not change dir to tmp: %v", err)
	}
	out := &bytes.Buffer{}
	app := &App{
		Out:     out,
		Printer: msg.Default(),
		Config: &config.Config{
			Tag: config.Tag{
				Files: []config.File{
					{
						Path:    "README.md",
						Search:  "version = {{.Current}}",
						Replace: "version = {{.Next}}",
					},
				},
			},
		},
		Replace: true,
	}

	err = app.Patch(false, false, "Message")
	if err != nil {
		t.Fatalf("app.Patch returned an error: %v", err)
	}

	wantOut := "force: false\npush: false\nmessage: Message\n"
	if out.String() != wantOut {
		t.Errorf("Wrong stdout, got %s, wanted %s", out.String(), wantOut)
	}
}
