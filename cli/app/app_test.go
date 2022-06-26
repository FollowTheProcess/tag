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

// Set up a tmp dir
// populate it with some random files and folders
// git init, add, and commit
// Tag some stuff
// test the apps methods

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

	err = add.Run()
	if err != nil {
		t.Fatalf("Error adding files to test git repo: %v", err)
	}

	err = commit.Run()
	if err != nil {
		t.Fatalf("Error committing to the test git repo: %v", err)
	}

	err = firstTag.Run()
	if err != nil {
		t.Fatalf("Error issuing the first tag to test git repo: %v", err)
	}

	tearDown := func() { os.RemoveAll(tmp) }

	return tmp, tearDown
}

func TestAppPatch(t *testing.T) {
	tmp, tearDown := setup(t)
	defer tearDown()

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
