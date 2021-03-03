package calver

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// this is for testing purpose only
var Now = time.Now

type Version struct {
	pattern string
	micro   int
	time    time.Time
}

func (v *Version) Next() *Version {
	latest, _ := NewVersion(v.pattern, 0)

	if latest.CompareTo(v) > 0 {
		return latest
	}
	return &Version{
		pattern: v.pattern,
		micro:   v.micro + 1,
		time:    v.time,
	}
}

func (v *Version) String() string {
	result := v.pattern
	if v.micro > 1 && !strings.Contains(v.pattern, MICRO.representation) {
		result += ".MICRO"
	}
	for _, k := range strings.Split(result, ".") {
		convention := CONVENTIONS[k]
		val := convention.extract(v)
		result = strings.Replace(result, convention.representation, convention.Format(val), 1)
	}
	return result
}

func (v *Version) CompareTo(v2 *Version) int {
	if v.time.Unix() > v2.time.Unix() {
		return 1
	}
	if v.time.Unix() < v2.time.Unix() {
		return -1
	}
	if v.micro > v2.micro {
		return 1
	}
	if v.micro < v2.micro {
		return -1
	}
	return 0
}

func ValidatePattern(pattern string) error {
	for _, segment := range strings.Split(pattern, ".") {
		validSegment := false
		for k := range CONVENTIONS {
			if segment == k {
				validSegment = true
				break
			}
		}
		if !validSegment {
			return fmt.Errorf("invalid pattern %s", pattern)
		}
	}

	return nil
}

func patternContainsSegment(patternSegments []string, segment string) bool {
	for _, s := range patternSegments {
		if s == segment {
			return true
		}
	}
	return false
}

func Parse(pattern, value string) (*Version, error) {
	var year, month, day, micro int

	if err := ValidatePattern(pattern); err != nil {
		return nil, err
	}

	patternSegments := strings.Split(pattern, ".")
	valueSegments := strings.Split(value, ".")
	numSegmentsPattern := len(patternSegments)
	numSegmentsValue := len(valueSegments)

	if numSegmentsValue-numSegmentsPattern == 1 &&
		!patternContainsSegment(patternSegments, MICRO.representation) {
		// if pattern doesn't contains MICRO segment but value does it should be considered valid and parsed
		patternSegments = append(patternSegments, MICRO.representation)
	} else if numSegmentsPattern != numSegmentsValue {
		return nil, fmt.Errorf("number of segments on value does not match pattern")
	}

	for i, pat := range patternSegments {
		val, err := strconv.Atoi(valueSegments[i])
		if err != nil {
			return nil, fmt.Errorf("invalid value")
		}
		segment := CONVENTIONS[pat]
		if err := segment.validate(val); err != nil {
			return nil, err
		}

		switch pat {
		case YYYY.representation, YY.representation, zeroY.representation:
			year = val
		case MM.representation, M0.representation, zeroM.representation:
			month = val
		case DD.representation, D0.representation, zeroD.representation:
			day = val
		case MICRO.representation:
			micro = val
		}
	}

	now := today()
	if year == 0 {
		year = now.Year()
	}
	if month == 0 {
		month = int(now.Month())
	}
	if day == 0 {
		day = now.Day()
	}
	if micro == 0 {
		micro = 1
	}
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return &Version{
		pattern: pattern,
		micro:   micro,
		time:    date,
	}, nil
}

func today() time.Time {
	y, m, d := Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func NewVersion(pattern string, micro int) (*Version, error) {
	if err := ValidatePattern(pattern); err != nil {
		return nil, err
	}

	if micro == 0 {
		micro = 1
	}
	return &Version{
		pattern: pattern,
		micro:   micro,
		time:    today(),
	}, nil
}
