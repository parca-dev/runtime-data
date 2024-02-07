package version

import "github.com/Masterminds/semver/v3"

// MustParseConstraints is a helper that wraps a call to NewConstraint and panics if the error is non-nil.
func MustParseConstraints(s string) *semver.Constraints {
	c, err := semver.NewConstraint(s)
	if err != nil {
		panic(err)
	}
	return c
}
