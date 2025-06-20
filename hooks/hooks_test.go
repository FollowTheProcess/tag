package hooks_test

import (
	"bytes"
	"testing"

	"go.followtheprocess.codes/tag/hooks"
)

func TestRunOk(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	if err := hooks.Run(hooks.StagePreCommit, "echo hello there", stdout, stderr); err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	want := "hello there\n"

	if got := stdout.String(); got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestRunNotOk(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	if err := hooks.Run(hooks.StagePreCommit, "exit 1", stdout, stderr); err == nil {
		t.Fatal("Run did not return an error")
	}
}

func TestRunNoOp(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	if err := hooks.Run(hooks.StagePreReplace, "", stdout, stderr); err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stdout.String() != "" {
		t.Errorf("Expected no stdout, got %s", stdout.String())
	}

	if stderr.String() != "" {
		t.Errorf("Expected no stderr, got %s", stderr.String())
	}
}
