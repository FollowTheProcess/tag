package version_test

import (
	"reflect"
	"testing"

	"github.com/FollowTheProcess/tag/version"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    version.Version
		wantErr bool
	}{
		{
			name:    "simple",
			text:    "1.2.4",
			want:    version.Version{Major: 1, Minor: 2, Patch: 4},
			wantErr: false,
		},
		{
			name:    "simple with v",
			text:    "v1.2.4",
			want:    version.Version{Major: 1, Minor: 2, Patch: 4},
			wantErr: false,
		},
		{
			name:    "prerelease",
			text:    "v2.3.7-rc.1",
			want:    version.Version{Major: 2, Minor: 3, Patch: 7, Prerelease: "rc.1"},
			wantErr: false,
		},
		{
			name:    "prerelease and build",
			text:    "v8.1.0-rc.1+build.123",
			want:    version.Version{Major: 8, Minor: 1, Patch: 0, Prerelease: "rc.1", Buildmetadata: "build.123"},
			wantErr: false,
		},
		{
			name:    "beta",
			text:    "1.2.3-beta",
			want:    version.Version{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta"},
			wantErr: false,
		},
		{
			name:    "obviously wrong",
			text:    "moby dick",
			want:    version.Version{},
			wantErr: true,
		},
		{
			name:    "invalid",
			text:    "1",
			want:    version.Version{},
			wantErr: true,
		},
		{
			name:    "prerelease digits",
			text:    "1.2.3-0123",
			want:    version.Version{},
			wantErr: true,
		},
		{
			name:    "extra parts",
			text:    "1.2.3.4",
			want:    version.Version{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := version.Parse(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() returned %v, wanted %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, wanted %#v", got, tt.want)
			}
		})
	}
}

func TestBumpMajor(t *testing.T) {
	tests := []struct {
		name    string
		current version.Version
		want    version.Version
	}{
		{
			name:    "zeros",
			current: version.Version{},
			want:    version.Version{Major: 1},
		},
		{
			name:    "minor",
			current: version.Version{Minor: 1},
			want:    version.Version{Major: 1},
		},
		{
			name:    "patch",
			current: version.Version{Patch: 1},
			want:    version.Version{Major: 1},
		},
		{
			name:    "everything",
			current: version.Version{Major: 0, Minor: 32, Patch: 6, Prerelease: "rc.1", Buildmetadata: "build.123"},
			want:    version.Version{Major: 1},
		},
		{
			name:    "big numbers",
			current: version.Version{Major: 123, Minor: 32, Patch: 6, Prerelease: "rc.1", Buildmetadata: "build.123"},
			want:    version.Version{Major: 124},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := version.BumpMajor(tt.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, wanted %#v", got, tt.want)
			}
		})
	}
}

func TestBumpMinor(t *testing.T) {
	tests := []struct {
		name    string
		current version.Version
		want    version.Version
	}{
		{
			name:    "zeros",
			current: version.Version{},
			want:    version.Version{Minor: 1},
		},
		{
			name:    "minor",
			current: version.Version{Minor: 1},
			want:    version.Version{Minor: 2},
		},
		{
			name:    "patch",
			current: version.Version{Patch: 1},
			want:    version.Version{Minor: 1},
		},
		{
			name:    "everything",
			current: version.Version{Major: 0, Minor: 32, Patch: 6, Prerelease: "rc.1", Buildmetadata: "build.123"},
			want:    version.Version{Minor: 33},
		},
		{
			name:    "big numbers",
			current: version.Version{Major: 123, Minor: 32, Patch: 6, Prerelease: "rc.1", Buildmetadata: "build.123"},
			want:    version.Version{Major: 123, Minor: 33},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := version.BumpMinor(tt.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, wanted %#v", got, tt.want)
			}
		})
	}
}

func TestBumpPatch(t *testing.T) {
	tests := []struct {
		name    string
		current version.Version
		want    version.Version
	}{
		{
			name:    "zeros",
			current: version.Version{},
			want:    version.Version{Patch: 1},
		},
		{
			name:    "minor",
			current: version.Version{Minor: 1},
			want:    version.Version{Minor: 1, Patch: 1},
		},
		{
			name:    "patch",
			current: version.Version{Patch: 1},
			want:    version.Version{Patch: 2},
		},
		{
			name:    "everything",
			current: version.Version{Major: 0, Minor: 32, Patch: 6, Prerelease: "rc.1", Buildmetadata: "build.123"},
			want:    version.Version{Minor: 32, Patch: 7},
		},
		{
			name:    "big numbers",
			current: version.Version{Major: 123, Minor: 32, Patch: 6, Prerelease: "rc.1", Buildmetadata: "build.123"},
			want:    version.Version{Major: 123, Minor: 32, Patch: 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := version.BumpPatch(tt.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, wanted %#v", got, tt.want)
			}
		})
	}
}
