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
