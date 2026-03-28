package contained_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/value/contained"
)

func TestTypeInfo_ValidateNonNilFuncs(t *testing.T) {
	require := require.New(t)
	require.NoError(contained.AddrPort().ValidateNonNilFuncs())
	require.NoError(contained.Bool().ValidateNonNilFuncs())
	require.NoError(contained.Duration().ValidateNonNilFuncs())
	require.NoError(contained.Float32().ValidateNonNilFuncs())
	require.NoError(contained.Float64().ValidateNonNilFuncs())
	require.NoError(contained.Int().ValidateNonNilFuncs())
	require.NoError(contained.Int16().ValidateNonNilFuncs())
	require.NoError(contained.Int32().ValidateNonNilFuncs())
	require.NoError(contained.Int64().ValidateNonNilFuncs())
	require.NoError(contained.Int8().ValidateNonNilFuncs())
	require.NoError(contained.NetIPAddr().ValidateNonNilFuncs())
	require.NoError(contained.Path().ValidateNonNilFuncs())
	require.NoError(contained.Rune().ValidateNonNilFuncs())
	require.NoError(contained.String().ValidateNonNilFuncs())
	require.NoError(contained.Uint().ValidateNonNilFuncs())
	require.NoError(contained.Uint16().ValidateNonNilFuncs())
	require.NoError(contained.Uint32().ValidateNonNilFuncs())
	require.NoError(contained.Uint64().ValidateNonNilFuncs())
	require.NoError(contained.Uint8().ValidateNonNilFuncs())
}
