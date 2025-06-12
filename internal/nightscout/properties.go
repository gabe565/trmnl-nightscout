package nightscout

import (
	"time"

	"gabe565.com/trmnl-nightscout/internal/config"
)

type Properties struct {
	Bgnow     Reading   `json:"bgnow"`
	Buckets   []Reading `json:"buckets"`
	Delta     Delta     `json:"delta"`
	Direction Direction `json:"direction"`
}

func (p Properties) String(conf *config.Config) string {
	result := p.Bgnow.DisplayBg(conf.Units) +
		" " + p.Bgnow.Arrow()
	if delta := p.Delta.Display(conf.Units); delta != "" {
		result += " " + p.Delta.Display(conf.Units)
	}
	if rel := p.Bgnow.Mills.Relative(true); rel != "" {
		result += " [" + p.Bgnow.Mills.Relative(true) + "]"
	}
	return result
}

func (p Properties) NextTimestamp() time.Time {
	if len(p.Buckets) == 0 {
		return time.Time{}
	}

	bucket := p.Buckets[0]
	lastDiff := bucket.ToMills.Sub(bucket.FromMills.Time)
	nextRead := p.Bgnow.Mills.Add(lastDiff)
	return nextRead
}

func (p Properties) Interval() time.Duration {
	mins, err := p.Delta.ElapsedMins.Float64()
	if err != nil {
		return 5 * time.Minute
	}
	return time.Duration(mins * float64(time.Minute)).Round(time.Minute)
}

func (p Properties) IsRecent(delay time.Duration) bool {
	interval := p.Interval()
	diff := time.Since(p.Bgnow.Mills.Add(delay))
	return diff <= interval
}
