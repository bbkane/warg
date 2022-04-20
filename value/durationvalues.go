package value

import "time"

func durationFromIFace(iFace interface{}) (time.Duration, error) {
	under, ok := iFace.(string)
	if !ok {
		return 0, ErrIncompatibleInterface
	}
	return durationFromString(under)
}

func durationFromString(s string) (time.Duration, error) {
	decoded, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return decoded, nil
}

// Duration is updateable from a string parsed with https://pkg.go.dev/time#ParseDuration
func Duration() (Value, error) {
	s := newScalarValue(
		0,
		"duration",
		fromIFaceFunc[time.Duration](durationFromIFace),
		fromStringFunc[time.Duration](durationFromString),
	)
	return &s, nil
}
