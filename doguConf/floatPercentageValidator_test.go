package doguConf

import (
	"strconv"
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/require"
)

func TestFloatPercentageValidator_CheckValidator_check(t *testing.T) {
	tests := []struct {
		name     string
		valid    []string
		input    string
		expected float64
		wantErr  bool
	}{
		// wrong format
		{name: "invalid input for negative value -200.00", input: "-200.00", wantErr: true},
		{name: "invalid input for negative value -200", input: "-200", wantErr: true},
		{name: "invalid input 0 missing decimals", input: "0", wantErr: true},
		{name: "invalid input 1000.00 too many preceding digits", input: "1000.0", wantErr: true},
		{name: "invalid input 70.555 too many succeeding digits", input: "70.555", wantErr: true},
		// min/max
		{name: "valid min 0.00", input: "0.00", expected: 0.00, wantErr: false},
		{name: "invalid min -0.01", input: "-0.01", wantErr: true},
		{name: "valid max 100.00", input: "100.00", expected: 100.00, wantErr: false},
		{name: "invalid max 100.01", input: "100.01", wantErr: true},
		// valid inputs
		{name: "valid input 10.00", input: "10.00", expected: 10.00, wantErr: false},
		{name: "valid input 15.35", input: "15.35", expected: 15.35, wantErr: false},
		{name: "valid input 99.99", input: "99.99", expected: 99.99, wantErr: false},
		{name: "valid input 53.0", input: "53.0", expected: 53.0, wantErr: false},
		{name: "valid input 74.45", input: "74.45", expected: 74.45, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := CreateFloatPercentageValidator(core.ValidationDescriptor{Type: FloatPercentageHundredKey})
			if err := o.Check(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				value, err := strconv.ParseFloat(tt.input, 64)
				require.NoError(t, err)
				require.Equal(t, tt.expected, value)
			}
		})
	}

	t.Run("error should contain min/max description", func(t *testing.T) {
		sut := CreateFloatPercentageValidator(core.ValidationDescriptor{Type: FloatPercentageHundredKey})

		err := sut.Check("150.00")

		require.Error(t, err)
		require.Contains(t, err.Error(), "should be between 0 and 100; invalid percentage")
	})
}
