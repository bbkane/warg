// package path provides a simple wrapper around a string path that can expand the users home directory, a common CLI need. Might be extracted into its own library and removed from warg if I find myself needing it in other code.
package path

import "github.com/mitchellh/go-homedir"

type Path struct {
	path string
}

func New(path string) Path {
	return Path{path: path}
}

func (p Path) String() string {
	return p.path
}

func (p Path) Expand() (string, error) {
	// TODO: rm homedir dependency by copying https://github.com/mitchellh/go-homedir/blob/af06845cf3004701891bf4fdb884bfe4920b3727/homedir.go#L58 and using os.UserHomeDir()
	expanded, err := homedir.Expand(p.path)
	if err != nil {
		return "", err
	}
	return expanded, nil
}

func (p Path) MustExpand() string {
	expanded, err := p.Expand()
	if err != nil {
		panic(err)
	}
	return expanded
}
