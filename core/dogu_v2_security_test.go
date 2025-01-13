package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

func Test_validateSecurity(t *testing.T) {
	type args struct {
		dogu *Dogu
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"valid empty", args{&Dogu{}}, assert.NoError},
		{"valid add filled", args{&Dogu{Security: Security{Capabilities: Capabilities{Add: []Capability{AuditControl}}}}}, assert.NoError},
		{"valid add filled", args{&Dogu{Security: Security{Capabilities: Capabilities{Drop: []Capability{AuditControl}}}}}, assert.NoError},
		{"all possible values", args{&Dogu{Security: Security{Capabilities: Capabilities{Add: AllCapabilities, Drop: AllCapabilities}}}}, assert.NoError},
		{"add all keyword", args{&Dogu{Security: Security{Capabilities: Capabilities{Add: []Capability{All}}}}}, assert.NoError},
		{"drop all keyword", args{&Dogu{Security: Security{Capabilities: Capabilities{Drop: []Capability{All}}}}}, assert.NoError},

		{"invalid valid add filled", args{&Dogu{Security: Security{Capabilities: Capabilities{Add: []Capability{"err"}}}}}, assert.Error},
		{"invalid valid drop filled", args{&Dogu{Security: Security{Capabilities: Capabilities{Drop: []Capability{"err"}}}}}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, tt.args.dogu.ValidateSecurity(), fmt.Sprintf("validateSecurity(%v)", tt.args.dogu))
		})
	}
}

func Test_validateSecurity_message(t *testing.T) {
	t.Run("should match for drop errors", func(t *testing.T) {
		// given
		dogu := &Dogu{Name: "official/dogu", Version: "1.2.3", Security: Security{Capabilities: Capabilities{Drop: []Capability{"err"}}}}

		// when
		actual := dogu.ValidateSecurity()

		// then
		require.Error(t, actual)
		assert.ErrorContains(t, actual, "dogu descriptor official/dogu:1.2.3 contains at least one invalid security field: err is not a valid capability to be dropped")
	})
	t.Run("should match for add errors", func(t *testing.T) {
		// given
		dogu := &Dogu{Name: "official/dogu", Version: "1.2.3", Security: Security{Capabilities: Capabilities{Add: []Capability{"err"}}}}

		// when
		actual := dogu.ValidateSecurity()

		// then
		require.Error(t, actual)
		assert.ErrorContains(t, actual, "dogu descriptor official/dogu:1.2.3 contains at least one invalid security field: err is not a valid capability to be added")
	})
}

func TestDogu_EffectiveCapabilities(t *testing.T) {
	type fields struct {
		Security Security
	}
	tests := []struct {
		name   string
		fields fields
		want   []Capability
	}{
		{"drop all", fields{Security{Capabilities: Capabilities{Drop: []Capability{All}}}}, []Capability{}},
		{"add all", fields{Security{Capabilities: Capabilities{Add: []Capability{All}}}}, AllCapabilities},
		{"drop all, add all", fields{Security{Capabilities: Capabilities{Drop: []Capability{All}, Add: []Capability{All}}}}, AllCapabilities},
		{"default list", fields{Security{Capabilities: Capabilities{}}}, DefaultCapabilities},
		{"drop every cap without all keyword", fields{Security{Capabilities: Capabilities{Drop: DefaultCapabilities}}}, nil},
		{"add 1 new and 1 existing caps to default list", fields{Security{Capabilities: Capabilities{Add: []Capability{Bpf, Chown}}}}, joinCapability(DefaultCapabilities, Bpf, Chown)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dogu{
				Security: tt.fields.Security,
			}
			assert.ElementsMatch(t, tt.want, d.EffectiveCapabilities(), "ListCapabilities()")
		})
	}
}

func joinCapability(capSlice []Capability, singleCaps ...Capability) []Capability {
	result := []Capability{}
	result = append(result, capSlice...)
	for _, singleCap := range singleCaps {
		if slices.Contains(capSlice, singleCap) {
			continue
		}

		result = append(result, singleCap)
	}
	return result
}
