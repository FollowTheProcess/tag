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
