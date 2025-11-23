package contained

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuneFromIFace(t *testing.T) {
	tests := []struct {
		name         string
		iFace        interface{}
		expectedRune rune
		expectedErr  bool
	}{
		{
			name:         "rune",
			iFace:        'a',
			expectedRune: 'a',
			expectedErr:  false,
		},
	}
	runeInner := Rune()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRune, actualErr := runeInner.FromIFace(tt.iFace)

			if tt.expectedErr {
				require.NotNil(t, actualErr)
			} else {
				require.Nil(t, actualErr)
			}

			require.Equal(t, tt.expectedRune, actualRune)
		})
	}

}
