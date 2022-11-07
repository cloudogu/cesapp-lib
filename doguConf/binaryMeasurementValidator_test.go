package doguConf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/core"
)

func TestBinaryMeasurementValidator_check(t *testing.T) {
	tests := []struct {
		name    string
		valid   []string
		input   string
		wantErr bool
	}{
		// min/max
		{name: "invalid input for negative value -200m", input: "-200m", wantErr: true},
		{name: "valid input 0b", input: "0b", wantErr: false},
		{name: "valid input max int64", input: "9223372036854775807g", wantErr: false},
		{name: "invalid input more than 19 digits", input: "10000000000000000000g", wantErr: true},
		// unit variation
		{name: "valid input 2b", input: "2b", wantErr: false},
		{name: "valid input 2k", input: "2k", wantErr: false},
		{name: "valid input 2m", input: "2m", wantErr: false},
		{name: "valid input 2g", input: "2g", wantErr: false},
		{name: "valid input 1024m", input: "1024m", wantErr: false},
		// human configuration error
		{name: "invalid input without unit", input: "2048", wantErr: true},
		{name: "invalid missing input value only units", input: "g", wantErr: true},
		{name: "invalid space between value and unit", input: "2 b", wantErr: true},
		{name: "invalid input uppercase unit 1B", input: "1B", wantErr: true},
		{name: "invalid input uppercase unit 1K", input: "1K", wantErr: true},
		{name: "invalid input uppercase unit 1M", input: "1M", wantErr: true},
		{name: "invalid input uppercase unit 1G", input: "1G", wantErr: true},
		{name: "invalid digit group separator ,", input: "1,024m", wantErr: true},
		{name: "invalid digit group separator _", input: "1_024m", wantErr: true},
		{name: "invalid digit group separator '", input: `1'024m`, wantErr: true},
		{name: "invalid decimal separator .", input: "1.024m", wantErr: true},
		// weird input
		{name: "invalid input for switched position", input: "m2", wantErr: true},
		{name: "invalid empty input", input: "", wantErr: true},
		{name: "invalid space input", input: " ", wantErr: true},
		{name: "invalid tab input", input: "	", wantErr: true},
		{name: "invalid input with unit x", input: "2x", wantErr: true},
		{name: "invalid characters input after unit", input: "1234masdfghj", wantErr: true},
		{name: "invalid unites all at once", input: "1234bkmg", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := CreateBinaryMeasurementValidator(core.ValidationDescriptor{Type: BinaryMeasurementKey})
			if err := o.Check(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("error should contain allowed binary measurement units", func(t *testing.T) {
		sut := CreateBinaryMeasurementValidator(core.ValidationDescriptor{Type: BinaryMeasurementKey})

		err := sut.Check("")

		require.Error(t, err)
		require.Contains(t, err.Error(), "valid units are: b,k,m,g")
	})
}

func TestBinaryMeasurementValidator_getUnitlessValue(t *testing.T) {
	t.Run("should return digit part without unit", func(t *testing.T) {
		sutInterfaced := CreateBinaryMeasurementValidator(core.ValidationDescriptor{Type: BinaryMeasurementKey})
		sut := sutInterfaced.(*BinaryMeasurementValidator)
		assert.NoError(t, sut.Check("123m"))

		actualValue, actualUnit := sut.SplitValueAndUnit("123m")

		require.Equal(t, "123", actualValue)
		require.Equal(t, "m", actualUnit)
	})
}
