package calver

import (
	"fmt"
	"math"
)

type Convention struct {
	representation string
	format         string
	extract        func(*Version) int
	validate       func(int) error
}

func (c *Convention) Format(value int) string {
	return fmt.Sprintf(c.format, value)
}

var (
	YYYY = Convention{
		representation: "YYYY",
		extract:        extractYear,
		format:         "%04d",
		validate:       validateInRange(2000, 2500),
	}
	YY = Convention{
		representation: "YY",
		extract:        truncInt(2, extractYear),
		format:         "%d",
		validate:       validateInRange(0, 99),
	}
	zeroY = Convention{
		representation: "0Y",
		extract:        truncInt(2, extractYear),
		format:         "%02d",
		validate:       validateInRange(0, 99),
	}
	MM = Convention{
		representation: "MM",
		extract:        extractMonth,
		format:         "%d",
		validate:       validateInRange(1, 12),
	}
	M0 = Convention{
		representation: "M0",
		extract:        extractMonth,
		format:         "%02d",
		validate:       validateInRange(1, 12),
	}
	zeroM = Convention{
		representation: "0M",
		extract:        extractMonth,
		format:         "%02d",
		validate:       validateInRange(1, 12),
	}
	DD = Convention{
		representation: "DD",
		extract:        extractDay,
		format:         "%d",
		validate:       validateInRange(1, 31),
	}
	D0 = Convention{
		representation: "D0",
		extract:        extractDay,
		format:         "%02d",
		validate:       validateInRange(1, 31),
	}
	zeroD = Convention{
		representation: "0D",
		extract:        extractDay,
		format:         "%02d",
		validate:       validateInRange(1, 31),
	}
	MICRO = Convention{
		representation: "MICRO",
		extract:        extractMicro,
		format:         "%d",
		validate:       validatePositive,
	}
	CONVENTIONS = map[string]Convention{
		YYYY.representation:  YYYY,
		YY.representation:    YY,
		zeroY.representation: zeroY,
		MM.representation:    MM,
		M0.representation:    M0,
		zeroM.representation: zeroM,
		DD.representation:    DD,
		D0.representation:    D0,
		zeroD.representation: zeroD,
		MICRO.representation: MICRO,
	}
)

func extractYear(version *Version) int {
	return version.time.Year()
}

func extractMonth(version *Version) int {
	return int(version.time.Month())
}

func extractDay(version *Version) int {
	return version.time.Day()
}

func extractMicro(version *Version) int {
	return version.micro
}

func truncInt(width int, fn func(*Version) int) func(*Version) int {
	return func(version *Version) int {
		multiplier := int(math.Pow10(width))
		val := fn(version)
		truncated := val - (val / multiplier * multiplier)
		return truncated
	}
}

func validateInRange(minVal, maxVal int) func(int) error {
	return func(val int) error {
		if minVal > val || maxVal < val {
			return fmt.Errorf("invalid value %d", val)
		}
		return nil
	}
}

func validatePositive(val int) error {
	if val <= 0 {
		return fmt.Errorf("invalid value %d", val)
	}
	return nil
}
