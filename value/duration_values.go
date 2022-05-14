package value

import (
	"time"

	str2duration "github.com/xhit/go-str2duration/v2"
)

func durationFromIFace(iFace interface{}) (time.Duration, error) {
	under, ok := iFace.(string)
	if !ok {
		return 0, ErrIncompatibleInterface
	}
	return durationFromString(under)
}

func durationFromString(s string) (time.Duration, error) {
	decoded, err := str2duration.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return decoded, nil
}

// Duration is updateable from a string parsed with https://pkg.go.dev/github.com/xhit/go-str2duration/v2#ParseDuration . Examples: "30s" (30 seconds), "3d" (3 days)
func Duration() (Value, error) {
	s := newScalarValue(
		0,
		"duration",
		fromIFaceFunc[time.Duration](durationFromIFace),
		fromStringFunc[time.Duration](durationFromString),
	)
	return &s, nil
}
