package value

import (
	"testing"

	"github.com/alecthomas/assert"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRune, actualErr := runeFromIFace(tt.iFace)

			if tt.expectedErr {
				assert.NotNil(t, actualErr)
			} else {
				assert.Nil(t, actualErr)
			}

			assert.Equal(t, tt.expectedRune, actualRune)
		})
	}

}
