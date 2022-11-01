package keys_test

import (
	"github.com/cloudogu/cesapp-lib/keys"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"

	"github.com/stretchr/testify/assert"
)

var emptyKeyProvider = ""

func TestInitRegister(t *testing.T) {
	assert.Contains(t, keys.KeyProviders, "pkcs1v15")
}

func TestEncryptDecrypt(t *testing.T) {
	provider, err := keys.NewKeyProvider(emptyKeyProvider)
	assert.Nil(t, err)
	pair, err := provider.Generate()
	assert.Nil(t, err)
	enc, err := pair.Public().Encrypt("hello cesapp")
	assert.Nil(t, err)
	assert.NotEqual(t, enc, "hello cesapp")
	dec, err := pair.Private().Decrypt(enc)
	assert.Nil(t, err)
	assert.Equal(t, "hello cesapp", dec)
}

func TestEncryptDecryptPayload(t *testing.T) {
	reg := registry.MockRegistry{}

	provider, err := keys.NewKeyProvider(emptyKeyProvider)
	assert.Nil(t, err)

	pair, err := provider.Generate()
	assert.Nil(t, err)

	pub, _ := pair.Public().AsString()
	err = reg.DoguConfig("myDogu").Set(registry.KeyDoguPublicKey, pub)
	assert.Nil(t, err)

	dogu := core.Dogu{
		Name: "unofficial/myDogu",
	}

	publicKeyString, err := reg.DoguConfig(dogu.GetSimpleName()).Get(registry.KeyDoguPublicKey)
	assert.NoError(t, err)

	publicKey, err := provider.ReadPublicKeyFromString(publicKeyString)
	assert.NoError(t, err)

	enc, err := publicKey.Encrypt("myPayload")
	assert.Nil(t, err)
	assert.NotEqual(t, enc, "myPayload")

	dec, err := pair.Private().Decrypt(enc)
	assert.Nil(t, err)
	assert.Equal(t, "myPayload", dec)
}

func TestEncryptDecryptLong(t *testing.T) {
	provider, err := keys.NewKeyProvider(emptyKeyProvider)
	assert.Nil(t, err)
	pair, err := provider.Generate()
	assert.Nil(t, err)
	textWith65Chars := "WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW"
	enc, err := pair.Public().Encrypt(textWith65Chars)
	assert.Nil(t, err)
	assert.NotEqual(t, enc, textWith65Chars)
	assert.True(t, strings.HasPrefix(enc, "{"))
	dec, err := pair.Private().Decrypt(enc)
	assert.Nil(t, err)
	assert.Equal(t, textWith65Chars, dec)
}

func TestCreateKeyPairFromPrivateKey(t *testing.T) {
	provider, err := keys.NewKeyProvider(emptyKeyProvider)
	assert.Nil(t, err)

	pair, err := provider.Generate()
	assert.Nil(t, err)
	enc, err := pair.Public().Encrypt("hello")
	assert.Nil(t, err)
	bytes, err := pair.Private().AsBytes()
	assert.Nil(t, err)

	pair, err = provider.FromPrivateKey(bytes)
	assert.Nil(t, err)
	dec, err := pair.Private().Decrypt(enc)
	assert.Nil(t, err)
	assert.Equal(t, "hello", dec)
}

func TestReadPublicKey(t *testing.T) {
	provider, err := keys.NewKeyProvider(emptyKeyProvider)
	assert.Nil(t, err)

	pair, err := provider.Generate()
	assert.Nil(t, err)
	puBytes, err := pair.Public().AsBytes()
	assert.Nil(t, err)

	publicKey, err := provider.ReadPublicKey(puBytes)
	assert.Nil(t, err)

	enc, err := publicKey.Encrypt("hello again")
	assert.Nil(t, err)
	dec, err := pair.Private().Decrypt(enc)
	assert.Nil(t, err)
	assert.Equal(t, "hello again", dec)
}

func TestPrivateKeyAsString(t *testing.T) {
	provider, err := keys.NewKeyProvider(emptyKeyProvider)
	assert.Nil(t, err)

	pair, err := provider.Generate()
	assert.Nil(t, err)

	enc, err := pair.Public().Encrypt("hello enc")
	assert.Nil(t, err)

	pk, err := pair.Private().AsString()
	assert.Nil(t, err)

	pair, err = provider.FromPrivateKey([]byte(pk))
	assert.Nil(t, err)
	dec, err := pair.Private().Decrypt(enc)
	assert.Nil(t, err)
	assert.Equal(t, "hello enc", dec)
}

func TestPublicKeyAsString(t *testing.T) {
	provider, err := keys.NewKeyProvider(emptyKeyProvider)
	assert.Nil(t, err)

	pair, err := provider.Generate()
	assert.Nil(t, err)

	pk, err := pair.Public().AsString()
	assert.Nil(t, err)

	publicKey, err := provider.ReadPublicKey([]byte(pk))
	assert.Nil(t, err)

	enc, err := publicKey.Encrypt("hello again")
	assert.Nil(t, err)
	dec, err := pair.Private().Decrypt(enc)
	assert.Nil(t, err)
	assert.Equal(t, "hello again", dec)
}

func TestKeyPairFromPrivateKeyPath(t *testing.T) {
	provider, err := keys.NewKeyProvider(emptyKeyProvider)
	assert.Nil(t, err)

	pair, err := provider.Generate()
	assert.Nil(t, err)

	enc, err := pair.Public().Encrypt("hello enc")
	assert.Nil(t, err)

	f, err := ioutil.TempFile("", "privatekey.pem")
	assert.Nil(t, err)
	f.Close()
	defer os.Remove(f.Name())

	err = pair.Private().ToFile(f.Name())
	assert.Nil(t, err)

	pair, err = provider.FromPrivateKeyPath(f.Name())
	assert.Nil(t, err)

	dec, err := pair.Private().Decrypt(enc)
	assert.Nil(t, err)
	assert.Equal(t, "hello enc", dec)
}
