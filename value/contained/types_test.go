package contained_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/value/contained"
)

func TestTypeInfo_ValidateNonNilFuncs(t *testing.T) {
	require := require.New(t)
	require.NoError(contained.AddrPort().ValidateNonNilFuncs())
	require.NoError(contained.Bool().ValidateNonNilFuncs())
	require.NoError(contained.Duration().ValidateNonNilFuncs())
	require.NoError(contained.DateTimeRFC3339().ValidateNonNilFuncs())
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

func TestDateTimeRFC3339(t *testing.T) {
	require := require.New(t)
	typeInfo := contained.DateTimeRFC3339()

	expected, err := time.Parse(time.RFC3339, "2026-01-02T03:04:05Z")
	require.NoError(err)

	fromString, err := typeInfo.FromString("2026-01-02T03:04:05Z")
	require.NoError(err)
	require.True(typeInfo.Equals(expected, fromString))

	fromIFace, err := typeInfo.FromIFace("2026-01-02T03:04:05Z")
	require.NoError(err)
	require.True(typeInfo.Equals(expected, fromIFace))

	_, err = typeInfo.FromString("not-a-date")
	require.Error(err)
}
