package doguConf

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
)

func TestOneOfValidator_check(t *testing.T) {
	tests := []struct {
		name    string
		valid   []string
		input   string
		wantErr bool
	}{
		{name: "valid value abc", valid: []string{"abc", "def"}, input: "abc", wantErr: false},
		{name: "valid value def", valid: []string{"abc", "def"}, input: "def", wantErr: false},
		{name: "valid value empty string", valid: []string{""}, input: "", wantErr: false},
		{name: "invalid value cde", valid: []string{"abc", "def"}, input: "cde", wantErr: true},
		{name: "invalid value abcdef", valid: []string{"abc", "def"}, input: "abcdef", wantErr: true},
		{name: "invalid value empty list", valid: []string{}, input: "abc", wantErr: true},
		{name: "invalid value empty input ", valid: []string{"abc", "def"}, input: "", wantErr: true},
		{name: "invalid value empty both ", valid: []string{}, input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := createOneOfValidator(core.ValidationDescriptor{Type: "ONE_OF", Values: tt.valid})
			if err := o.Check(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
