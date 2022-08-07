// Package version implements tag's semantic version parsing and bumping functionality.
package version

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	major         = "major"
	minor         = "minor"
	patch         = "patch"
	prerelease    = "prerelease"
	buildmetadata = "buildmetadata"
)

// See https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
var semVerRegex = regexp.MustCompile(fmt.Sprintf(`^v?(?P<%s>0|[1-9]\d*)\.(?P<%s>0|[1-9]\d*)\.(?P<%s>0|[1-9]\d*)(?:-(?P<%s>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<%s>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`, major, minor, patch, prerelease, buildmetadata))

// Version encodes a semantic version.
type Version struct {
	Prerelease    string
	Buildmetadata string
	Major         int
	Minor         int
	Patch         int
}

// String satisfies the stringer interface and allows a Version to print itself.
//
//	v := Version{Major: 1, Minor: 2, Patch: 3}
//	fmt.Println(v) // "1.2.3"
func (v Version) String() string {
	base := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		base += "-" + v.Prerelease
	}
	if v.Buildmetadata != "" {
		base += "+" + v.Buildmetadata
	}
	return base
}

// Tag creates a string representation of the version suitable for git tags
// it is identical to the String() method except prepends a 'v' to the result.
//
//	v := Version{Major: 1, Minor: 2, Patch: 3}
//	fmt.Println(v) // "v1.2.3"
func (v Version) Tag() string {
	return "v" + v.String()
}

// Parse creates a Version from a semver string.
func Parse(text string) (Version, error) {
	if !semVerRegex.MatchString(text) {
		return Version{}, fmt.Errorf("%q is not a valid semantic version", text)
	}

	groups := make(map[string]string, 5) // 5 elements (fields of Version)

	parts := semVerRegex.FindStringSubmatch(text)
	names := semVerRegex.SubexpNames()

	for i, part := range parts {
		// First element in parts is the substring e.g. "v1.2.4"
		// first in names is empty string -> skip both
		if i == 0 {
			continue
		}
		groups[names[i]] = part
	}

	majorInt, err := strconv.Atoi(groups[major])
	if err != nil {
		return Version{}, fmt.Errorf("major part %v cannot be cast to int", groups[major])
	}
	minorInt, err := strconv.Atoi(groups[minor])
	if err != nil {
		return Version{}, fmt.Errorf("minor part %v cannot be cast to int", groups[minor])
	}
	patchInt, err := strconv.Atoi(groups[patch])
	if err != nil {
		return Version{}, fmt.Errorf("patch part %v cannot be cast to int", groups[patch])
	}

	v := Version{
		Prerelease:    groups[prerelease],
		Buildmetadata: groups[buildmetadata],
		Major:         majorInt,
		Minor:         minorInt,
		Patch:         patchInt,
	}

	return v, nil
}

// BumpMajor returns a new Version with it's major version bumped.
func BumpMajor(current Version) Version {
	// Everything else set to zero value
	return Version{
		Major: current.Major + 1,
	}
}

// BumpMinor returns a new Version with it's minor version bumped.
func BumpMinor(current Version) Version {
	// Keep major, bump minor, everything else -> zero value
	return Version{
		Major: current.Major,
		Minor: current.Minor + 1,
	}
}

// BumpPatch returns a new Version with it's patch version bumped.
func BumpPatch(current Version) Version {
	// Keep major and minor, bump patch, everything else -> zero value
	return Version{
		Major: current.Major,
		Minor: current.Minor,
		Patch: current.Patch + 1,
	}
}
