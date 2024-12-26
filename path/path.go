// package path provides a simple wrapper around a string path that can expand the users home directory, a common CLI need. Might be extracted into its own library and removed from warg if I find myself needing it in other code.
package path

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Path struct {
	path string
}

func New(path string) Path {
	return Path{path: path}
}

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
		return "", errors.New("cannot expand user-specific home dir")
	}

	return filepath.Join(homedir, p.path[1:]), nil
}

// Expand expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
func (p Path) Expand() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home dir: %w", err)
	}
	expanded, err := p.expand(homedir)
	if err != nil {
		return "", err
	}
	return expanded, nil
}

// MustExpand calls `Expand` and panics on any errors
func (p Path) MustExpand() string {
	expanded, err := p.Expand()
	if err != nil {
		panic(err)
	}
	return expanded
}
