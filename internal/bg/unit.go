//nolint:recvcheck
package bg

const MmolConversionFactor = 0.0555

//go:generate go tool enumer -type Unit -linecomment -text

type Unit uint8

const (
	Mgdl Unit = iota // mg/dL
	Mmol             // mmol/L
)
