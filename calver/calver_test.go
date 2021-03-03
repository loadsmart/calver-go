package calver

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockNowFunc(fn func() time.Time) func() {
	Now = fn
	return func() {
		Now = time.Now
	}
}

func TestMain(m *testing.M) {
	reset := mockNowFunc(func() time.Time {
		return time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)
	})
	defer reset()
	os.Exit(m.Run())
}

func TestNewVersion_Format(t *testing.T) {
	var dateTests = []struct {
		name    string
		format  string
		initial string
	}{
		{
			name:    "YYYY.MM.DD format",
			format:  "YYYY.MM.DD",
			initial: "2021.2.1",
		},
		{
			name:    "YY.MM.DD format",
			format:  "YY.MM.DD",
			initial: "21.2.1",
		},
		{
			name:    "YY.0M.0D format",
			format:  "YY.0M.0D",
			initial: "21.02.01",
		},
		{
			name:    "0Y.MM.DD format",
			format:  "0Y.MM.DD",
			initial: "21.2.1",
		},
		{
			name:    "0Y.0M.0D format",
			format:  "0Y.0M.0D",
			initial: "21.02.01",
		},
		{
			name:    "YY.MM format",
			format:  "YY.MM",
			initial: "21.2",
		},
	}

	for _, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewVersion(tt.format, 0)
			finalVersion := tt.initial
			assert.NoError(t, err)
			for i := 0; i <= 10; i++ {
				if i != 0 {
					finalVersion = fmt.Sprintf("%s.%d", tt.initial, i+1)
					c = c.Next()
				}
				assert.Exactly(t, finalVersion, c.String())
			}
		})
	}
}

func TestParse(t *testing.T) {
	var dateTests = []struct {
		name    string
		format  string
		initial string
		error   bool
	}{
		{
			name:    "YYYY.MM.DD format",
			format:  "YYYY.MM.DD",
			initial: "2007.1.1",
			error:   false,
		},
		{
			name:    "YYYY.MM.DD.MICRO format",
			format:  "YYYY.MM.DD.MICRO",
			initial: "2007.1.1.1",
			error:   false,
		},
		{
			name:    "wrong format",
			format:  "YYYY.MM.DD",
			initial: "2007.1.1dev",
			error:   true,
		},
		{
			name:    "wrong value",
			format:  "YYYY.MM.DD",
			initial: "2007.1.1-dev.99",
			error:   true,
		},
		{
			name:    "bigger format",
			format:  "YYYY.MM.MICRO",
			initial: "2007.1000",
			error:   true,
		},
	}
	for _, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(tt.format, tt.initial)
			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Exactly(t, c.String(), tt.initial)
			}
		})
	}
}

func TestParse_Release(t *testing.T) {
	reset := mockNowFunc(func() time.Time {
		return time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)
	})
	defer reset()

	var dateTests = []struct {
		name    string
		format  string
		initial string
		final   string
	}{
		{
			name:    "YYYY.MM.DD.MICRO format",
			format:  "YYYY.MM.DD.MICRO",
			initial: "2021.2.1.1",
			final:   "2021.2.1.2",
		},
		{
			name:    "YYYY.MM.DD format",
			format:  "YYYY.MM.DD",
			initial: "2021.2.1",
			final:   "2021.2.1.2",
		},
		{
			name:    "YYYY.MM format",
			format:  "YYYY.MM",
			initial: "2021.2",
			final:   "2021.2.2",
		},
		{
			name:    "YYYY.MM.MICRO format",
			format:  "YYYY.MM.MICRO",
			initial: "2021.2.4",
			final:   "2021.2.5",
		},
		{
			name:    "YYYY.DD.MICRO format",
			format:  "YYYY.DD.MICRO",
			initial: "2021.1.1",
			final:   "2021.1.2",
		},
		{
			name:    "0M.0D.YYYY.MICRO format",
			format:  "0M.0D.YYYY.MICRO",
			initial: "02.01.2021.1",
			final:   "02.01.2021.2",
		},
	}
	for _, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(tt.format, tt.initial)
			assert.NoError(t, err)
			assert.Exactly(t, tt.initial, c.String())
			assert.Exactly(t, tt.final, c.Next().String())
		})
	}
}

func TestNew_DifferentDay(t *testing.T) {
	n := 0
	reset := mockNowFunc(func() time.Time {
		return time.Date(2007, 2, 1+n, 0, 0, 0, 0, time.UTC)
	})
	defer reset()

	c, _ := NewVersion("YYYY.MM.DD", 0)

	for _, d := range []int{1, 2, 3, 4, 5} {
		date := fmt.Sprintf("2007.2.%d", d)

		t.Run(date, func(t *testing.T) {
			assert.Exactly(t, date, c.String())
			n = d
			c = c.Next()
		})
	}
}

func TestNew_DifferentMonth(t *testing.T) {
	m := time.January
	reset := mockNowFunc(func() time.Time {
		return time.Date(2007, m, 5, 0, 0, 0, 0, time.UTC)
	})
	defer reset()

	var dateTests = []struct {
		name  string
		month time.Month
	}{
		{
			name:  "month january",
			month: time.January,
		},
		{
			name:  "month february",
			month: time.February,
		},
		{
			name:  "month march",
			month: time.March,
		},
		{
			name:  "month april",
			month: time.April,
		},
		{
			name:  "month may",
			month: time.May,
		},
		{
			name:  "month june",
			month: time.June,
		},
	}
	for i, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			m = tt.month
			c, _ := NewVersion("YYYY.MM.DD", 0)
			finalVersion := fmt.Sprintf("2007.%d.5", i+1)
			assert.Exactly(t, finalVersion, c.String())
		})
	}
}

func TestNew_DifferentYear(t *testing.T) {
	n := 0
	reset := mockNowFunc(func() time.Time {
		return time.Date(2001+n, 2, 5, 0, 0, 0, 0, time.UTC)
	})
	defer reset()

	c, _ := NewVersion("YYYY.MM.DD", 0)

	for _, d := range []int{1, 2, 3, 4, 5} {
		date := fmt.Sprintf("200%d.2.5", d)

		t.Run(date, func(t *testing.T) {
			assert.Exactly(t, date, c.String())
			n = d
			c = c.Next()
		})
	}
}

func TestNew_UnsupportedFormat(t *testing.T) {
	_, err := NewVersion("YYYY.XX.HH", 0)
	assert.EqualError(t, err, "invalid pattern YYYY.XX.HH")
}

func TestVersion_CompareTo(t *testing.T) {

	var tests = []struct {
		name string
		v1       Version
		v2       Version
		expected int
	}{
		{
			name: "equal",
			v1: Version{
				time:  time.Date(2108, 5, 9, 0, 0, 0, 0, time.UTC),
				micro: 1,
			},
			v2: Version{
				time:  time.Date(2108, 5, 9, 0, 0, 0, 0, time.UTC),
				micro: 1,
			},
			expected: 0,
		},
		{
			name: "later date",
			v1: Version{
				time:  time.Date(2108, 5, 10, 0, 0, 0, 0, time.UTC),
				micro: 1,
			},
			v2: Version{
				time:  time.Date(2108, 5, 9, 0, 0, 0, 0, time.UTC),
				micro: 1,
			},
			expected: 1,
		},
		{
			name: "earlier date",
			v1: Version{
				time:  time.Date(2108, 5, 8, 0, 0, 0, 0, time.UTC),
				micro: 1,
			},
			v2: Version{
				time:  time.Date(2108, 5, 9, 0, 0, 0, 0, time.UTC),
				micro: 1,
			},
			expected: -1,
		},
		{
			name: "bigger micro",
			v1: Version{
				time:  time.Date(2108, 5, 9, 0, 0, 0, 0, time.UTC),
				micro: 2,
			},
			v2: Version{
				time:  time.Date(2108, 5, 9, 0, 0, 0, 0, time.UTC),
				micro: 1,
			},
			expected: 1,
		},
		{
			name: "smaller micro",
			v1: Version{
				time:  time.Date(2108, 5, 9, 0, 0, 0, 0, time.UTC),
				micro: 1,
			},
			v2: Version{
				time:  time.Date(2108, 5, 9, 0, 0, 0, 0, time.UTC),
				micro: 2,
			},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.v1.CompareTo(&tt.v2))
		})
	}
}
