package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/FollowTheProcess/tag/config"
)

func TestLoad(t *testing.T) {
	path := makeConfigFile(t)

	got, err := config.Load(path)
	if err != nil {
		t.Fatalf("config.Load returned an error: %v", err)
	}

	want := &config.Config{
		Tag: config.Tag{
			Files: []config.File{
				{
					Path:    "hello.go",
					Search:  "version {{.Current}}",
					Replace: "version {{.Next}}",
				},
				{
					Path:    "another.go",
					Search:  "version {{.Current}}",
					Replace: "version {{.Next}}",
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, wanted %#v", got, want)
	}
}

func TestLoadAlt(t *testing.T) {
	path := makeAltConfigFile(t)

	got, err := config.Load(path)
	if err != nil {
		t.Fatalf("config.Load returned an error: %v", err)
	}

	want := &config.Config{
		Tag: config.Tag{
			Files: []config.File{
				{
					Path:    "hello.go",
					Search:  "version {{.Current}}",
					Replace: "version {{.Next}}",
				},
				{
					Path:    "another.go",
					Search:  "version {{.Current}}",
					Replace: "version {{.Next}}",
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, wanted %#v", got, want)
	}
}

func TestRender(t *testing.T) {
	conf := &config.Config{
		Tag: config.Tag{
			Files: []config.File{
				{
					Path:    "hello.go",
					Search:  "version {{.Current}}",
					Replace: "version {{.Next}}",
				},
				{
					Path:    "another.go",
					Search:  "version {{.Current}}",
					Replace: "version {{.Next}}",
				},
			},
		},
	}

	if err := conf.Render("0.1.0", "0.2.0"); err != nil {
		t.Fatalf("Render returned an error: %v", err)
	}

	want := []config.File{
		{
			Path:    "hello.go",
			Search:  "version 0.1.0",
			Replace: "version 0.2.0",
		},
		{
			Path:    "another.go",
			Search:  "version 0.1.0",
			Replace: "version 0.2.0",
		},
	}

	if !reflect.DeepEqual(conf.Files, want) {
		t.Errorf("got %#v, wanted %#v", conf.Files, want)
	}
}

// makeConfigFile creates a temporary config file, returning it's path.
func makeConfigFile(t *testing.T) string {
	file, err := os.CreateTemp("", "tag.toml")
	if err != nil {
		t.Fatalf("CreateTemp returned an error: %v", err)
	}
	defer file.Close()
	doc := `
	[tag]
	files = [
		{ path = "hello.go", search = "version {{.Current}}", replace = "version {{.Next}}"},
		{ path = "another.go", search = "version {{.Current}}", replace = "version {{.Next}}"},
	]`
	_, err = file.WriteString(doc)
	if err != nil {
		t.Fatalf("Could not write to tmp file: %v", err)
	}

	return file.Name()
}

// makeAltConfigFile creates a temporary config file with alternative TOML syntax
// returning it's path.
func makeAltConfigFile(t *testing.T) string {
	file, err := os.CreateTemp("", "tag.toml")
	if err != nil {
		t.Fatalf("CreateTemp returned an error: %v", err)
	}
	defer file.Close()
	doc := `
	[[tag.files]]
	path = "hello.go"
	search = "version {{.Current}}"
	replace = "version {{.Next}}"

	[[tag.files]]
	path = "another.go"
	search = "version {{.Current}}"
	replace = "version {{.Next}}"
	`
	_, err = file.WriteString(doc)
	if err != nil {
		t.Fatalf("Could not write to tmp file: %v", err)
	}

	return file.Name()
}
