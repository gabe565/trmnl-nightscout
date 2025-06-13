package bg

type BG float64

func (m BG) Mgdl() float64 { return float64(m) }

func (m BG) Mmol() float64 { return float64(m) * MmolConversionFactor }

func (m BG) Value(u Unit) float64 {
	switch u {
	case Mgdl:
		return m.Mgdl()
	case Mmol:
		return m.Mmol()
	default:
		panic("invalid unit")
	}
}
