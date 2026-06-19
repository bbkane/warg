// Package path provides a file path wrapper that supports tilde (~) expansion
// to the user's home directory.
package path

import (
	"errors"
	"os"
	"path/filepath"

	"go.bbkane.com/warg/colerr"
)

// Path wraps a file system path string, providing tilde expansion via [Path.Expand].
type Path struct {
	path string
}

// New creates a [Path] from the given string. No validation is performed.
func New(path string) Path {
	return Path{path: path}
}

// String returns the raw path string without expansion.
func (p Path) String() string {
	return p.path
}

func (p Path) expand(homedir string) (string, error) {
	// adapted from https://github.com/mitchellh/go-homedir/blob/af06845cf3004701891bf4fdb884bfe4920b3727/homedir.go#L58
	if len(p.path) == 0 {
		return "", nil
	}

	if p.path[0] != '~' {
		return p.path, nil
	}

	if len(p.path) > 1 && p.path[1] != '/' && p.path[1] != '\\' {
		return "", errors.New("Cannot expand user-specific home dir")
	}

	return filepath.Join(homedir, p.path[1:]), nil
}

// Expand returns the path with a leading ~ replaced by the user's home directory.
// If the path does not start with ~, it is returned unchanged.
// Returns an error if the home directory cannot be determined.
func (p Path) Expand() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", colerr.NewWrapped(err, "Could not get home dir")
	}
	expanded, err := p.expand(homedir)
	if err != nil {
		return "", err
	}
	return expanded, nil
}

// MustExpand calls [Path.Expand] and panics on error.
func (p Path) MustExpand() string {
	expanded, err := p.Expand()
	if err != nil {
		panic(err)
	}
	return expanded
}

// Equals reports whether two paths have the same raw string (no expansion).
func (p Path) Equals(other Path) bool {
	return p.path == other.path
}
