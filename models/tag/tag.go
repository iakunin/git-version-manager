package tag

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"strings"
)

type BumpStrategy string

const (
	Patch BumpStrategy = "patch"
	Minor              = "minor"
	Major              = "major"
)

const prefixSuffixSeparator = "-"

type Tag struct {
	prefix  string
	suffix  string
	version *semver.Version
}

func Empty(prefix string, suffix string) *Tag {
	return &Tag{
		suffix:  suffix,
		prefix:  prefix,
		version: semver.New("0.0.0"),
	}
}

func New(rawTag string, prefix string, suffix string) (*Tag, error) {
	if prefix != "" && !strings.HasPrefix(rawTag, prefix) {
		return nil, errors.New(
			fmt.Sprintf("rawTag ('%s') must start with prefix '%s'", rawTag, prefix),
		)
	}

	if suffix != "" && !strings.HasSuffix(rawTag, suffix) {
		return nil, errors.New(
			fmt.Sprintf("rawTag ('%s') must ends with suffix '%s'", rawTag, suffix),
		)
	}

	version, err := parseVersion(rawTag, prefix, suffix)
	if err != nil {
		return nil, err
	}

	return &Tag{
		suffix:  suffix,
		prefix:  prefix,
		version: version,
	}, nil
}

func parseVersion(rawTag string, prefix string, suffix string) (*semver.Version, error) {
	version := rawTag

	if prefix != "" {
		version = strings.Replace(version, prefix+prefixSuffixSeparator, "", 1)
	}

	if suffix != "" {
		version = strings.Replace(version, prefixSuffixSeparator+suffix, "", 1)
	}

	newVersion, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	if newVersion.PreRelease != "" {
		return nil, errors.New(
			fmt.Sprintf("expected `PreRelease` to be empty, got: '%s'", newVersion.PreRelease),
		)
	}

	return newVersion, err
}

func (v *Tag) String() string {
	var buffer bytes.Buffer

	if v.prefix != "" {
		_, _ = fmt.Fprintf(&buffer, "%s%s", v.prefix, prefixSuffixSeparator)
	}

	_, _ = fmt.Fprintf(&buffer, "%s", v.version.String())

	if v.suffix != "" {
		_, _ = fmt.Fprintf(&buffer, "%s%s", prefixSuffixSeparator, v.suffix)
	}

	return buffer.String()
}

// Compare tests if v is less than, equal to, or greater than `another`,
// returning -1, 0, or +1 respectively.
func (v *Tag) Compare(another Tag) int {
	if cmp := strings.Compare(v.prefix, another.prefix); cmp != 0 {
		return cmp
	}

	if cmp := strings.Compare(v.suffix, another.suffix); cmp != 0 {
		return cmp
	}

	return v.version.Compare(*another.version)
}

// Equal tests if v is equal to `another`.
func (v *Tag) Equal(another Tag) bool {
	return v.Compare(another) == 0
}

// LessThan tests if v is less than `another`.
func (v *Tag) LessThan(another Tag) bool {
	return v.Compare(another) < 0
}

func (v *Tag) Bump(strategy BumpStrategy) {
	switch strategy {
	case Major:
		v.version.BumpMajor()
		return
	case Minor:
		v.version.BumpMinor()
		return
	case Patch:
		v.version.BumpPatch()
		return
	default:
		v.version.BumpPatch()
		return
	}
}
