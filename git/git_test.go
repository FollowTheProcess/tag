package git //nolint: testpackage // We need access to internals to mock os.Exec

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

var (
	mockExitStatus int
	mockStdout     string
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	es := strconv.Itoa(mockExitStatus)
	cmd.Env = []string{
		"GO_WANT_HELPER_PROCESS=1",
		"STDOUT=" + mockStdout,
		"EXIT_STATUS=" + es,
	}
	return cmd
}

func TestExecCommandHelper(t *testing.T) {
	t.Helper()
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	tmp := t.TempDir()
	t.Setenv("GOCOVERDIR", tmp)

	fmt.Fprint(os.Stdout, os.Getenv("STDOUT"))
	i, _ := strconv.Atoi(os.Getenv("EXIT_STATUS")) //nolint: errcheck // Ignore error here

	if err := os.RemoveAll(tmp); err != nil {
		t.Fatalf("could not remove tmp: %v", err)
	}
	os.Exit(i) //nolint: revive // Needed for the helper process
}

func TestCommit(t *testing.T) {
	tests := []struct {
		name    string
		stdout  string
		status  int
		wantErr bool
	}{
		{
			name:    "happy",
			stdout:  "success",
			status:  0,
			wantErr: false,
		},
		{
			name:    "sad",
			stdout:  "I failed!",
			status:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExitStatus = tt.status
			mockStdout = tt.stdout
			gitCommand = fakeExecCommand
			defer func() { gitCommand = exec.Command }()

			out, err := Commit("Bump version 0.1.0 -> 0.2.0")
			if (err != nil) != tt.wantErr {
				t.Fatalf("Commit() returned %v, wanted %v", err, tt.wantErr)
			}

			if out != tt.stdout {
				t.Errorf("Commit stdout was %q, wanted %q", out, tt.stdout)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		stdout  string
		status  int
		wantErr bool
	}{
		{
			name:    "happy",
			stdout:  "success",
			status:  0,
			wantErr: false,
		},
		{
			name:    "sad",
			stdout:  "I failed!",
			status:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExitStatus = tt.status
			mockStdout = tt.stdout
			gitCommand = fakeExecCommand
			defer func() { gitCommand = exec.Command }()

			err := Add()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Add() returned %v, wanted %v", err, tt.wantErr)
			}
		})
	}
}

func TestPush(t *testing.T) {
	tests := []struct {
		name    string
		stdout  string
		status  int
		wantErr bool
	}{
		{
			name:    "happy",
			stdout:  "success",
			status:  0,
			wantErr: false,
		},
		{
			name:    "sad",
			stdout:  "I failed!",
			status:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExitStatus = tt.status
			mockStdout = tt.stdout
			gitCommand = fakeExecCommand
			defer func() { gitCommand = exec.Command }()

			out, err := Push()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Push() returned %v, wanted %v", err, tt.wantErr)
			}

			if out != tt.stdout {
				t.Errorf("Push stdout was %q, wanted %q", out, tt.stdout)
			}
		})
	}
}

func TestListTags(t *testing.T) {
	tests := []struct {
		name    string
		stdout  string
		status  int
		wantErr bool
	}{
		{
			name: "happy",
			stdout: `
			v0.6.0
			v0.5.0
			v0.4.2
			v0.4.1
			v0.4.0
			v0.3.2
			v0.3.1
			v0.3.0
			v0.2.0`,
			status:  0,
			wantErr: false,
		},
		{
			name:    "sad",
			stdout:  "I failed!",
			status:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExitStatus = tt.status
			mockStdout = tt.stdout
			gitCommand = fakeExecCommand
			defer func() { gitCommand = exec.Command }()

			out, _, err := ListTags(10)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ListTags() returned %v, wanted %v", err, tt.wantErr)
			}

			if out != tt.stdout {
				t.Errorf("ListTags stdout was %q, wanted %q", out, tt.stdout)
			}
		})
	}
}

func TestLatestTag(t *testing.T) {
	tests := []struct {
		name    string
		stdout  string
		status  int
		wantErr bool
	}{
		{
			name:    "happy",
			stdout:  "v1.67.2",
			status:  0,
			wantErr: false,
		},
		{
			name:    "sad",
			stdout:  "I failed!",
			status:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExitStatus = tt.status
			mockStdout = tt.stdout
			gitCommand = fakeExecCommand
			defer func() { gitCommand = exec.Command }()

			out, err := LatestTag()
			if (err != nil) != tt.wantErr {
				t.Fatalf("LatestTag() returned %v, wanted %v", err, tt.wantErr)
			}

			if out != tt.stdout {
				t.Errorf("LatestTag stdout was %q, wanted %q", out, tt.stdout)
			}
		})
	}
}

func TestCreateTag(t *testing.T) {
	tests := []struct {
		name    string
		stdout  string
		status  int
		wantErr bool
	}{
		{
			name:    "happy",
			stdout:  "Woohoo tag created",
			status:  0,
			wantErr: false,
		},
		{
			name:    "sad",
			stdout:  "I failed!",
			status:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExitStatus = tt.status
			mockStdout = tt.stdout
			gitCommand = fakeExecCommand
			defer func() { gitCommand = exec.Command }()

			out, err := CreateTag("v1.4.5", "This is a tag")
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateTag() returned %v, wanted %v", err, tt.wantErr)
			}

			if out != tt.stdout {
				t.Errorf("CreateTag stdout was %q, wanted %q", out, tt.stdout)
			}
		})
	}
}

func TestIsRepo(t *testing.T) {
	tests := []struct {
		name   string
		stdout string
		status int
		want   bool
	}{
		{
			name:   "yes",
			stdout: "true",
			status: 0,
			want:   true,
		},
		{
			name:   "no",
			stdout: "fatal: not a git repository (or any of the parent directories): .git",
			status: 1,
			want:   false,
		},
		{
			name:   "some other error",
			stdout: "Argh!",
			status: 1,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExitStatus = tt.status
			mockStdout = tt.stdout
			gitCommand = fakeExecCommand
			defer func() { gitCommand = exec.Command }()

			if got := IsRepo(); got != tt.want {
				t.Errorf("IsRepo returned %v, wanted %v", got, tt.want)
			}
		})
	}
}

func TestIsDirty(t *testing.T) {
	tests := []struct {
		name    string
		stdout  string
		status  int
		want    bool
		wantErr bool
	}{
		{
			name: "yes",
			stdout: `
			M git/git.go
			M git_git_test.go
			`,
			status:  0,
			want:    true,
			wantErr: false,
		},
		{
			name:    "no",
			stdout:  "",
			status:  0,
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExitStatus = tt.status
			mockStdout = tt.stdout
			gitCommand = fakeExecCommand
			defer func() { gitCommand = exec.Command }()

			got, err := IsDirty()
			if (err != nil) != tt.wantErr {
				t.Fatalf("IsDirty() returned %v, wanted %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("IsDirty returned %v, wanted %v", got, tt.want)
			}
		})
	}
}

func TestBranch(t *testing.T) {
	tests := []struct {
		name    string
		stdout  string
		want    string
		status  int
		wantErr bool
	}{
		{
			name:    "main",
			stdout:  "main",
			status:  0,
			want:    "main",
			wantErr: false,
		},
		{
			name:    "master",
			stdout:  "master",
			status:  0,
			want:    "master",
			wantErr: false,
		},
		{
			name:    "trunk",
			stdout:  "trunk",
			status:  0,
			want:    "trunk",
			wantErr: false,
		},
		{
			name:    "bad",
			stdout:  "bad",
			status:  1,
			want:    "bad",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExitStatus = tt.status
			mockStdout = tt.stdout
			gitCommand = fakeExecCommand
			defer func() { gitCommand = exec.Command }()

			got, err := Branch()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Branch() returned %v, wanted %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("Branch() returned %s, wanted %s", got, tt.want)
			}
		})
	}
}
