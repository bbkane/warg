package path_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/path"
)

func TestPath(t *testing.T) {
	p := path.New("~/testpath")
	require.Equal(t, "~/testpath", p.String())
	// TODO: expand the path and test that
}
