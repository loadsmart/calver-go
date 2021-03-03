package calver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvention_Format(t *testing.T) {
	convention := Convention{
		format: "%05d",
	}
	formatted := convention.Format(123)

	assert.Equal(t, "00123", formatted)
}

func TestConventions(t *testing.T) {
	var tests = []struct {
		convention  Convention
		expectedValue  int
		expectedFormatted string
	}{
		{
			convention:        YYYY,
			expectedValue:     2108,
			expectedFormatted: "2108",
		},
		{
			convention:        YY,
			expectedValue:     8,
			expectedFormatted: "8",
		},
		{
			convention:        zeroY,
			expectedValue:     8,
			expectedFormatted: "08",
		},
		{
			convention:        MM,
			expectedValue:     5,
			expectedFormatted: "5",
		},
		{
			convention:        M0,
			expectedValue:     5,
			expectedFormatted: "05",
		},
		{
			convention:        zeroM,
			expectedValue:     5,
			expectedFormatted: "05",
		},
		{
			convention:        DD,
			expectedValue:     9,
			expectedFormatted: "9",
		},
		{
			convention:        D0,
			expectedValue:     9,
			expectedFormatted: "09",
		},
		{
			convention:        zeroD,
			expectedValue:     9,
			expectedFormatted: "09",
		},
		{
			convention:        MICRO,
			expectedValue:     15,
			expectedFormatted: "15",
		},
	}
	
	version := &Version{
		time: time.Date(2108, 5, 9, 0, 0, 0, 0, time.UTC),
		micro: 15,
	}

	for _, tt := range tests {
		t.Run(tt.convention.representation, func(t *testing.T) {
			value := tt.convention.extract(version)
			formatted := tt.convention.Format(value)
			assert.NoError(t, tt.convention.validate(value))
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedFormatted, formatted)
		})
	}
}

func TestValidateInRange(t *testing.T) {
	var tests = []struct {
		name        string
		value       int
		expectError bool
	}{
		{
			name:        "in range",
			value:       50,
			expectError: false,
		},
		{
			name:        "before range",
			value:       30,
			expectError: true,
		},
		{
			name:        "after range",
			value:       70,
			expectError: true,
		},
		{
			name:        "inside lower bound",
			value:       40,
			expectError: false,
		},
		{
			name:        "outside lower bound",
			value:       39,
			expectError: true,
		},
		{
			name:        "inside upper bound",
			value:       60,
			expectError: false,
		},
		{
			name:        "outside upper bound",
			value:       61,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := validateInRange(40, 60)
			err := fn(tt.value)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePositive(t *testing.T) {
	var tests = []struct {
		name        string
		value       int
		expectError bool
	}{
		{
			name:        "zero",
			value:       0,
			expectError: true,
		},
		{
			name:        "above 0",
			value:       1,
			expectError: false,
		},
		{
			name:        "below 0",
			value:       -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePositive(tt.value)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
