package json

import (
	stdjson "encoding/json"
	"time"
)

// Duration is exactly the same as time.Duration but (un)marshals using same format as https://pkg.go.dev/time#ParseDuration
type Duration struct {
	time.Duration
}

var _ stdjson.Marshaler = &Duration{}
var _ stdjson.Unmarshaler = &Duration{}

func (d *Duration) MarshalJSON() ([]byte, error) {
	s := d.String()
	return stdjson.Marshal(s)
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := stdjson.Unmarshal(b, &s); err != nil {
		return err
	}

	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = dur
	return nil
}
