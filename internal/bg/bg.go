package bg

func MgdlToMmol(v float64) float64 {
	return v * MmolConversionFactor
}

func MmolToMgdl(v float64) float64 {
	return v / MmolConversionFactor
}

func NewMgdl[T float64 | BG | int](v T) BG {
	return BG(v)
}

func NewMmol[T float64 | BG | int](v T) BG {
	return BG(MmolToMgdl(float64(v)))
}

type BG float64

func (m BG) Mgdl() float64 { return float64(m) }

func (m BG) Mmol() float64 { return MgdlToMmol(m.Mgdl()) }

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
