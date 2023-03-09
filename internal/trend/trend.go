package trend

type Trend string

const (
	Upward   Trend = "upward"
	Downward Trend = "downward"
	NoTrend  Trend = "none"
)

func Calculate(values []int) Trend {
	if len(values) == 0 {
		return NoTrend
	}

	if len(values) == 1 {
		return NoTrend
	}

	if values[0] > values[len(values)-1] {
		return Upward
	} else if values[0] < values[len(values)-1] {
		return Downward
	} else {
		return NoTrend
	}
}
