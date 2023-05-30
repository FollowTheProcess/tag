package config_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/FollowTheProcess/tag/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    config.Config
		wantErr bool
	}{
		{
			name: "fully populated",
			file: "full.toml",
			want: config.Config{
				Version: "0.1.0",
				Git: config.Git{
					DefaultBranch:   "master",
					MessageTemplate: "Custom version {{.Current}} -> {{.Next}}",
					TagTemplate:     "taggy v{{.Next}}",
				},
				Hooks: config.Hooks{
					PreReplace: "echo 'I run before doing anything'",
					PreCommit:  "echo 'I run after replacing but before committing changes'",
					PreTag:     "echo 'I run after committing changes but before tagging'",
					PrePush:    "echo 'I run after tagging, but before pushing'",
				},
				Files: []config.File{
					{
						Path:   "pyproject.toml",
						Search: `version = "{{.Current}}"`,
					},
					{
						Path:   "README.md",
						Search: "My project, version {{.Current}}",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing git",
			file: "nogit.toml",
			want: config.Config{
				Version: "0.1.0",
				Git: config.Git{
					DefaultBranch:   "main",
					MessageTemplate: "Bump version {{.Current}} -> {{.Next}}",
					TagTemplate:     "v{{.Next}}",
				},
				Hooks: config.Hooks{
					PreReplace: "echo 'I run before doing anything'",
					PreCommit:  "echo 'I run after replacing but before committing changes'",
					PreTag:     "echo 'I run after committing changes but before tagging'",
					PrePush:    "echo 'I run after tagging, but before pushing'",
				},
				Files: []config.File{
					{
						Path:   "pyproject.toml",
						Search: `version = "{{.Current}}"`,
					},
					{
						Path:   "README.md",
						Search: "My project, version {{.Current}}",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing hooks",
			file: "nohooks.toml",
			want: config.Config{
				Version: "0.1.0",
				Git: config.Git{
					DefaultBranch:   "main",
					MessageTemplate: "Bump version {{.Current}} -> {{.Next}}",
					TagTemplate:     "v{{.Next}}",
				},
				Hooks: config.Hooks{},
				Files: []config.File{
					{
						Path:   "pyproject.toml",
						Search: `version = "{{.Current}}"`,
					},
					{
						Path:   "README.md",
						Search: "My project, version {{.Current}}",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "minimal",
			file: "minimal.toml",
			want: config.Config{
				Version: "0.1.0",
				Git: config.Git{
					DefaultBranch:   "main",
					MessageTemplate: "Bump version {{.Current}} -> {{.Next}}",
					TagTemplate:     "v{{.Next}}",
				},
				Hooks: config.Hooks{},
				Files: nil,
			},
			wantErr: false,
		},
		{
			name:    "empty",
			file:    "empty.toml",
			want:    config.Config{},
			wantErr: true,
		},
		{
			name:    "not exists",
			file:    "missing.toml",
			want:    config.Config{},
			wantErr: true,
		},
		{
			name:    "not toml",
			file:    "nottoml.json",
			want:    config.Config{},
			wantErr: true,
		},
		{
			name:    "bad toml",
			file:    "bad.toml",
			want:    config.Config{},
			wantErr: true,
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get cwd: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(cwd, "testdata", tt.file)

			got, err := config.Load(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr = %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Got:\n%#v\n\nWanted:\n%#v\n", got, tt.want)
			}
		})
	}
}

func TestRender(t *testing.T) {
	cfg := config.Config{
		Version: "1.0.0",
		Git: config.Git{
			DefaultBranch:   "main",
			MessageTemplate: "Bump version {{.Current}} -> {{.Next}}",
			TagTemplate:     "v{{.Next}}",
		},
		Files: []config.File{
			{
				Path:   "some.txt",
				Search: "Version = {{.Current}}",
			},
			{
				Path:   "other.md",
				Search: "version = {{.Current}}",
			},
		},
	}
	if err := cfg.Render("1.0.0", "2.0.0"); err != nil {
		t.Fatalf("Render returned an error: %v", err)
	}

	want := config.Config{
		Version: "1.0.0",
		Git: config.Git{
			DefaultBranch:   "main",
			MessageTemplate: "Bump version 1.0.0 -> 2.0.0",
			TagTemplate:     "v2.0.0",
		},
		Files: []config.File{
			{
				Path:    "some.txt",
				Search:  "Version = 1.0.0",
				Replace: "Version = 2.0.0",
			},
			{
				Path:    "other.md",
				Search:  "version = 1.0.0",
				Replace: "version = 2.0.0",
			},
		},
	}

	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("Got:\n%#v\n\nWanted:\n%#v\n", cfg, want)
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		name string
		file string
		cfg  config.Config
	}{
		{
			name: "full",
			file: "full.toml",
			cfg: config.Config{
				Version: "0.1.0",
				Git: config.Git{
					DefaultBranch:   "master",
					MessageTemplate: "Custom version {{.Current}} -> {{.Next}}",
					TagTemplate:     "taggy v{{.Next}}",
				},
				Hooks: config.Hooks{
					PreReplace: "echo 'I run before doing anything'",
					PreCommit:  "echo 'I run after replacing but before committing changes'",
					PreTag:     "echo 'I run after committing changes but before tagging'",
					PrePush:    "echo 'I run after tagging, but before pushing'",
				},
				Files: []config.File{
					{
						Path:   "pyproject.toml",
						Search: `version = "{{.Current}}"`,
					},
					{
						Path:   "README.md",
						Search: "My project, version {{.Current}}",
					},
				},
			},
		},
		{
			name: "missing git",
			file: "nogit.toml",
			cfg: config.Config{
				Version: "0.1.0",
				Hooks: config.Hooks{
					PreReplace: "echo 'I run before doing anything'",
					PreCommit:  "echo 'I run after replacing but before committing changes'",
					PreTag:     "echo 'I run after committing changes but before tagging'",
					PrePush:    "echo 'I run after tagging, but before pushing'",
				},
				Files: []config.File{
					{
						Path:   "pyproject.toml",
						Search: `version = "{{.Current}}"`,
					},
					{
						Path:   "README.md",
						Search: "My project, version {{.Current}}",
					},
				},
			},
		},
		{
			name: "missing hooks",
			file: "nohooks.toml",
			cfg: config.Config{
				Version: "0.1.0",
				Files: []config.File{
					{
						Path:   "pyproject.toml",
						Search: `version = "{{.Current}}"`,
					},
					{
						Path:   "README.md",
						Search: "My project, version {{.Current}}",
					},
				},
			},
		},
		{
			name: "minimal",
			file: "minimal.toml",
			cfg: config.Config{
				Version: "0.1.0",
			},
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get cwd: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp, err := os.CreateTemp("", "config.toml*")
			if err != nil {
				t.Fatalf("could not create temp file: %v", err)
			}
			tmp.Close()
			defer os.RemoveAll(tmp.Name())
			if err = tt.cfg.Save(tmp.Name()); err != nil {
				t.Fatalf("Save returned an error: %v", err)
			}

			written, err := os.ReadFile(tmp.Name())
			if err != nil {
				t.Fatalf("could not read written file: %v", err)
			}

			written = bytes.ReplaceAll(written, []byte("\r\n"), []byte("\n"))

			golden, err := os.ReadFile(filepath.Join(cwd, "testdata", tt.file))
			if err != nil {
				t.Fatalf("could not read golden file: %v", err)
			}

			golden = bytes.ReplaceAll(golden, []byte("\r\n"), []byte("\n"))

			if string(written) != string(golden) {
				t.Errorf("Got:\n%#v\n\nWanted:\n%#v\n", string(written), string(golden))
			}
		})
	}
}
