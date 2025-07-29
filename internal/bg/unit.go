//nolint:recvcheck
package bg

import "bytes"

const MmolConversionFactor = 0.0555

//go:generate go tool enumer -type Unit -linecomment

type Unit uint8

const (
	Mgdl Unit = iota // mg/dL
	Mmol             // mmol/L
)

// MarshalText implements the encoding.TextMarshaler interface for Unit.
func (u Unit) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Unit.
func (u *Unit) UnmarshalText(text []byte) error {
	var err error
	if *u, err = UnitString(string(text)); err != nil {
		switch string(bytes.ToLower(text)) {
		case "mgdl":
			*u = Mgdl
			return nil
		case "mmol":
			*u = Mmol
			return nil
		default:
			return err
		}
	}
	return nil
}
