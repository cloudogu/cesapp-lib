package doguConf

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.etcd.io/etcd/client/v2"

	"github.com/cloudogu/cesapp-lib/core"
)

func TestConfigValidator_Check(t *testing.T) {
	type config struct {
		global map[string]string
		dogu   map[string]string
	}
	tests := []struct {
		name    string
		config  config
		field   core.ConfigurationField
		wantErr bool
	}{
		{name: "valid dogu entry", config: config{dogu: map[string]string{"key": "optionB"}},
			field:   core.ConfigurationField{Name: "key", Validation: core.ValidationDescriptor{Type: "ONE_OF", Values: []string{"optionA", "optionB"}}},
			wantErr: false},
		{name: "invalid dogu entry", config: config{dogu: map[string]string{"key": "optionC"}},
			field:   core.ConfigurationField{Name: "key", Validation: core.ValidationDescriptor{Type: "ONE_OF", Values: []string{"optionA", "optionB"}}},
			wantErr: true},
		{name: "valid global entry", config: config{global: map[string]string{"key": "optionB"}},
			field:   core.ConfigurationField{Name: "key", Global: true, Validation: core.ValidationDescriptor{Type: "ONE_OF", Values: []string{"optionA", "optionB"}}},
			wantErr: false},
		{name: "invalid global entry", config: config{global: map[string]string{"key": "optionC"}},
			field:   core.ConfigurationField{Name: "key", Global: true, Validation: core.ValidationDescriptor{Type: "ONE_OF", Values: []string{"optionA", "optionB"}}},
			wantErr: true},
		{name: "invalid encrypted entry", config: config{dogu: map[string]string{"key": "optionC"}},
			field:   core.ConfigurationField{Name: "key", Encrypted: true, Validation: core.ValidationDescriptor{Type: "ONE_OF", Values: []string{"optionA", "optionB"}}},
			wantErr: false},
		{name: "optional entry not found", config: config{dogu: map[string]string{"key": "notfound"}},
			field:   core.ConfigurationField{Name: "key", Optional: true, Validation: core.ValidationDescriptor{Type: "ONE_OF", Values: []string{"optionA", "optionB"}}},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			reader := newMockConfigReader(t)
			for key, value := range tt.config.dogu {
				if value == "notfound" {
					reader.On("Get", key).Return("", client.Error{Code: client.ErrorCodeKeyNotFound})
				} else {
					reader.On("Get", key).Return(value, nil)
				}
			}
			for key, value := range tt.config.global {
				if value == "notfound" {
					reader.On("GetGlobal", key).Return(value, client.Error{Code: client.ErrorCodeKeyNotFound})
				} else {
					reader.On("GetGlobal", key).Return(value, nil)
				}
			}

			c := &ConfigValidator{
				configReader: reader,
			}

			// when then
			if err := c.Check(tt.field); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
			reader.AssertExpectations(t)
		})
	}
}

func TestConfigValidator_CheckWithoutType(t *testing.T) {
	reader := newMockConfigReader(t)
	name := "key"
	reader.On("Get", name).Return("", nil)

	validator := &ConfigValidator{configReader: reader}
	err := validator.Check(core.ConfigurationField{
		Name:       name,
		Validation: core.ValidationDescriptor{},
	})

	assert.Nil(t, err)
}

func TestConfigValidator_CheckErrorFromReader(t *testing.T) {
	reader := newMockConfigReader(t)
	name := "key"
	expectedError := client.Error{}
	reader.
		On("Get", name).
		Return("", expectedError)

	validator := &ConfigValidator{configReader: reader}
	err := validator.Check(core.ConfigurationField{
		Name:       name,
		Validation: core.ValidationDescriptor{},
	})

	assert.Error(t, err, expectedError)
}
