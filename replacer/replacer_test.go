package replacer

import (
	"fmt"
	"os"
	"testing"
)

func TestReplace(t *testing.T) {
	t.Run("search term exists", func(t *testing.T) {
		file, err := os.CreateTemp("", "test.txt")
		if err != nil {
			t.Fatalf("CreateTemp returned an error: %v", err)
		}

		_, err = file.WriteString("Hello I'm some file contents, this should say: REPLACE ME")
		if err != nil {
			t.Fatalf("Could not write to tmp file: %v", err)
		}

		file.Close()

		if err = Replace(file.Name(), "REPLACE ME", "NEW STRING"); err != nil {
			t.Fatalf("Replace returned an error: %v", err)
		}

		newContents, err := os.ReadFile(file.Name())
		if err != nil {
			t.Fatalf("Could not read from replaced file: %v", err)
		}

		want := "Hello I'm some file contents, this should say: NEW STRING"

		if got := string(newContents); got != want {
			t.Errorf("Wrong file contents: got %q, wanted %q", got, want)
		}
	})

	t.Run("search term doesnt exist", func(t *testing.T) {
		file, err := os.CreateTemp("", "test.txt")
		if err != nil {
			t.Fatalf("CreateTemp returned an error: %v", err)
		}

		_, err = file.WriteString("Hello I'm some file contents, this should say: I'm not here")
		if err != nil {
			t.Fatalf("Could not write to tmp file: %v", err)
		}

		file.Close()

		if err = Replace(file.Name(), "REPLACE ME", "NEW STRING"); err == nil {
			t.Fatal("Replace should have returned an error but got nil")
		}

		want := fmt.Sprintf("Could not find %q in %s", "REPLACE ME", file.Name())

		if err.Error() != want {
			t.Errorf("Wrong error: got %q, wanted %q", err.Error(), want)
		}

		newContents, err := os.ReadFile(file.Name())
		if err != nil {
			t.Fatalf("Could not read from replaced file: %v", err)
		}

		want = "Hello I'm some file contents, this should say: I'm not here"

		if got := string(newContents); got != want {
			t.Errorf("Wrong file contents: got %q, wanted %q", got, want)
		}
	})
}
