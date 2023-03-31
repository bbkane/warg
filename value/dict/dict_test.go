package dict_test

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
	"go.bbkane.com/warg/value/dict"
)

func TestDict_New(t *testing.T) {
	constructor := dict.New(contained.Int(), dict.Default(map[string]int{"one": 1}), dict.Choices(1, 2))
	v, err := constructor()
	require.Nil(t, err)
	dictVal := v.(value.DictValue)

	require.Equal(t, map[string]string{"one": "1"}, dictVal.DefaultStringMap())
	require.Equal(t, []string{"1", "2"}, dictVal.Choices())
}

func TestDict_ReplaceFromInterface(t *testing.T) {

	tests := []struct {
		name          string
		update        interface{}
		expectedValue interface{}
		expectedErr   error
	}{
		{
			name: "fine",
			update: map[string]interface{}{
				"hi": 1,
			},
			expectedValue: map[string]int{
				"hi": 1,
			},
			expectedErr: nil,
		},
		{

			name: "badKey",
			update: map[int]interface{}{
				1: 1,
			},
			expectedValue: map[string]int{},
			expectedErr:   contained.ErrIncompatibleInterface,
		},
		{
			name: "badVal",
			update: map[string]interface{}{
				"one": "bad",
			},
			expectedValue: map[string]int{},
			expectedErr:   contained.ErrIncompatibleInterface,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constructor := dict.New(contained.Int())
			v, err := constructor()
			require.Nil(t, err)
			dictVal := v.(value.DictValue)

			actualErr := dictVal.ReplaceFromInterface(tt.update)
			require.Equal(t, tt.expectedErr, actualErr)
			require.Equal(t, tt.expectedValue, dictVal.Get())
		})
	}
}

func TestDict_Update(t *testing.T) {
	constructor := dict.New(contained.Addr())
	v, err := constructor()
	require.Nil(t, err)
	dictVal := v.(value.DictValue)

	err = dictVal.Update("key=1.1.1.1")
	require.Nil(t, err)
	expected := map[string]netip.Addr{
		"key": netip.MustParseAddr("1.1.1.1"),
	}
	actual := dictVal.Get().(map[string]netip.Addr)
	require.Equal(t, expected, actual)
}
