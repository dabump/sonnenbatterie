package trend

type Trend string

const (
	Upward   Trend = "upward"
	Downward Trend = "downward"
	NoTrend  Trend = "none"
)

func Calculate(values []float64) Trend {
	if len(values) == 0 || len(values) == 1 {
		return NoTrend
	}

	if values[0] > values[len(values)-1] {
		return Upward
	} else if values[0] < values[len(values)-1] {
		return Downward
	}

	return NoTrend
}
