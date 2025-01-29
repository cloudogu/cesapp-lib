package doguConf

import (
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
)

func TestHasConfiguration(t *testing.T) {
	assert.False(t, HasConfiguration(&core.Dogu{}))
	assert.True(t, HasConfiguration(&core.Dogu{Configuration: []core.ConfigurationField{
		{
			Name: "title",
		},
	}}))
}

func TestEditConfigurationShouldOverwriteValueForKey1ButKeepValueForKey2(t *testing.T) {
	// given
	configKeyTitle := "title"
	configKeyLimit := "limit"
	config1 := core.ConfigurationField{
		Name:        configKeyTitle,
		Description: "my desc",
		Default:     "defaultVal",
		Validation: core.ValidationDescriptor{
			Type:   OneOfKey,
			Values: []string{"a", "b", "value-from-reader"},
		},
	}
	config2 := core.ConfigurationField{
		Name:        configKeyLimit,
		Description: "a limit in bytes with lowercase ",
		Optional:    true,
		Validation: core.ValidationDescriptor{
			Type: BinaryMeasurementKey,
		},
	}

	writer := NewMockFieldWriter(t)
	writer.EXPECT().Print(config1, "hitchhiker to the galaxy").Return(nil).Once()
	writer.EXPECT().Print(config2, "hitchhiker to the galaxy").Return(nil).Once()
	reader := NewMockFieldReader(t)
	reader.EXPECT().Read().Return("value-from-reader", nil).Once()
	reader.EXPECT().Read().Return("1024m", nil).Once()

	configurationContext := newMockDoguConfigurationContext(t)
	configurationContextExpect := configurationContext.EXPECT()
	configurationContextExpect.Exists(configKeyTitle).Return(true, nil).Once()
	configurationContextExpect.Get(configKeyTitle).Return("hitchhiker to the galaxy", nil).Once()
	configurationContextExpect.Set(configKeyTitle, "value-from-reader").Return(nil).Once()
	configurationContextExpect.Exists(configKeyLimit).Return(true, nil).Once()
	configurationContextExpect.Get(configKeyLimit).Return("hitchhiker to the galaxy", nil).Once()
	configurationContextExpect.Set(configKeyLimit, "1024m").Return(nil).Once()
	editor := DoguConfigurationEditor{
		ConfigurationContext: configurationContext,
		Writer:               writer,
		Reader:               reader,
	}

	// when
	err := editor.EditConfiguration([]core.ConfigurationField{config1, config2}, false)

	// then
	require.NoError(t, err)
}

func TestEditConfigurationOfAEncryptedField(t *testing.T) {
	// given
	writer := NewMockFieldWriter(t)
	field := core.ConfigurationField{Name: "title", Encrypted: true}
	writer.EXPECT().Print(field, "_encrypted_").Return(nil)
	reader := NewMockFieldReader(t)
	reader.EXPECT().Read().Return("input", nil)

	keyProvider, err := keys.NewKeyProvider("pkcs1v15")
	assert.Nil(t, err)
	keyPair, err := keyProvider.Generate()
	assert.Nil(t, err)

	configurationContext := newMockDoguConfigurationContext(t)
	configurationExpect := configurationContext.EXPECT()
	configurationExpect.Exists("title").Return(true, nil)
	configurationExpect.Set("title", mocks.Anything).Return(nil)
	editor := DoguConfigurationEditor{
		ConfigurationContext: configurationContext,
		Writer:               writer,
		Reader:               reader,
		PublicKey:            keyPair.Public(),
	}

	// when
	err = editor.EditConfiguration([]core.ConfigurationField{field}, false)

	// then
	assert.Nil(t, err)
}

func TestEditConfigurationWithEmptyValue(t *testing.T) {
	// given
	writer := NewMockFieldWriter(t)
	field := core.ConfigurationField{Name: "title"}
	writer.EXPECT().Print(field, "hitchhiker").Return(nil)
	reader := NewMockFieldReader(t)
	reader.EXPECT().Read().Return("input", nil)

	configurationContext := newMockDoguConfigurationContext(t)
	configurationExpect := configurationContext.EXPECT()
	configurationExpect.Exists("title").Return(true, nil)
	configurationExpect.Get("title").Return("hitchhiker", nil)
	configurationExpect.Set("title", "input").Return(nil)
	editor := DoguConfigurationEditor{
		ConfigurationContext: configurationContext,
		Writer:               writer,
		Reader:               reader,
	}

	// when

	err := editor.EditConfiguration([]core.ConfigurationField{field}, false)

	// then
	assert.Nil(t, err)
}

func TestEditConfigurationOfAGlobalField(t *testing.T) {
	// given
	writer := NewMockFieldWriter(t)
	field := core.ConfigurationField{Name: "title", Global: true}
	reader := NewMockFieldReader(t)

	configurationContext := newMockDoguConfigurationContext(t)
	editor := DoguConfigurationEditor{
		ConfigurationContext: configurationContext,
		Writer:               writer,
		Reader:               reader,
	}

	// when
	err := editor.EditConfiguration([]core.ConfigurationField{field}, false)

	// then
	assert.Nil(t, err)
}

func TestEditConfigurationOfADirectory(t *testing.T) {
	// given
	writer := NewMockFieldWriter(t)
	field := core.ConfigurationField{Name: "title", IsDirectory: true}
	reader := NewMockFieldReader(t)

	configurationContext := newMockDoguConfigurationContext(t)
	editor := DoguConfigurationEditor{
		ConfigurationContext: configurationContext,
		Writer:               writer,
		Reader:               reader,
	}

	// when
	err := editor.EditConfiguration([]core.ConfigurationField{field}, false)

	// then
	assert.Nil(t, err)
}

func TestEditConfigurationWithDeleteOnEmpty(t *testing.T) {
	// given
	writer := NewMockFieldWriter(t)
	field := core.ConfigurationField{Name: "title"}
	writer.EXPECT().Print(field, "hitchhiker").Return(nil)
	reader := NewMockFieldReader(t)
	reader.EXPECT().Read().Return("", nil)

	configurationContext := newMockDoguConfigurationContext(t)
	configurationExpect := configurationContext.EXPECT()
	configurationExpect.Exists("title").Return(true, nil)
	configurationExpect.Get("title").Return("hitchhiker", nil)
	configurationExpect.Delete("title").Return(nil)
	editor := DoguConfigurationEditor{
		ConfigurationContext: configurationContext,
		Writer:               writer,
		Reader:               reader,
	}

	// when
	err := editor.EditConfiguration([]core.ConfigurationField{field}, true)

	// then
	assert.Nil(t, err)
}

func TestEditConfigurationWithDeleteOnEmptyOnANonExistingKey(t *testing.T) {
	// given
	writer := NewMockFieldWriter(t)
	field := core.ConfigurationField{Name: "title"}
	writer.EXPECT().Print(field, "").Return(nil)
	reader := NewMockFieldReader(t)
	reader.EXPECT().Read().Return("", nil)

	configurationContext := newMockDoguConfigurationContext(t)
	configurationExpect := configurationContext.EXPECT()
	configurationExpect.Exists("title").Return(false, nil)
	editor := DoguConfigurationEditor{
		ConfigurationContext: configurationContext,
		Writer:               writer,
		Reader:               reader,
	}

	// when
	err := editor.EditConfiguration([]core.ConfigurationField{field}, true)

	// then
	assert.Nil(t, err)
}
