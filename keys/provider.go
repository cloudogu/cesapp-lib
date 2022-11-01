package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

const (
	bitSize int = 2048
)

// Encrypter encrypts a given reader stream with a given public key.
// This method may be exported into a library and must not be unexported.
type Encrypter func(random io.Reader, pub *rsa.PublicKey, msg []byte) ([]byte, error)

// Decrypter decrypts a given reader stream with a given public key.
// This method may be exported into a library and must not be unexported.
type Decrypter func(rand io.Reader, priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error)

// KeyProvider provides functions for en- and decryption.
// This method may be exported into a library and must not be unexported.
type KeyProvider struct {
	Encrypter Encrypter
	Decrypter Decrypter
}

// NewKeyProvider creates a new KeyProvider.
// This method may be exported into a library and must not be unexported.
func NewKeyProvider(keyType string) (*KeyProvider, error) {
	provider := providers[keyType]
	if provider == nil {
		return nil, errors.New("could not find provider from type " + keyType)
	}
	return provider, nil
}

// Generate creates a new public/private key
func (provider *KeyProvider) Generate() (*KeyPair, error) {
	log.Info("create new key pair")
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rsa key pair")
	}
	return &KeyPair{key: key, encrypter: provider.Encrypter, decrypter: provider.Decrypter}, nil
}

// FromPrivateKeyPath reads the keypair from the private key file path
func (provider *KeyProvider) FromPrivateKeyPath(path string) (*KeyPair, error) {
	pk, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read private key file")
	}
	return provider.FromPrivateKey(pk)
}

// FromPrivateKey creates a key pair from the private key
func (provider *KeyProvider) FromPrivateKey(privateKey []byte) (*KeyPair, error) {
	p, _ := pem.Decode(privateKey)
	key, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private rsa key")
	}

	return &KeyPair{key: key, encrypter: provider.Encrypter, decrypter: provider.Decrypter}, nil
}

// ReadPublicKeyFromString reads a public key from its string representation
func (provider *KeyProvider) ReadPublicKeyFromString(publicKeyString string) (*PublicKey, error) {
	return provider.ReadPublicKey([]byte(publicKeyString))
}

// ReadPublicKey reads a public key from an byte array
func (provider *KeyProvider) ReadPublicKey(publicKey []byte) (*PublicKey, error) {
	p, _ := pem.Decode(publicKey)
	key, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse public rsa key")
	}
	pk, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("could not cast key to rsa.PublicKey")
	}

	return &PublicKey{pk, provider.Encrypter}, nil
}
