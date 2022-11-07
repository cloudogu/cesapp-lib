package keys

import (
	"crypto/rsa"
	"encoding/base64"
)

// AesKeyBitLength is the used aes key length
const AesKeyBitLength = 256

// NonceBitLength is the used nonce length
const NonceBitLength = 96

// MaxRSAEncryptionLength defines the max length of values which will be encrypted with RSA instead of hybrid encryption
const MaxRSAEncryptionLength = 64

// AesGcm represents AES with the block cipher mode GCM
const AesGcm = "AES_GCM"

// HybridEncryptionValue contains an encrypted value and information about the encryption
type HybridEncryptionValue struct {
	Encryption Encryption `json:"encryption"`
	Value      string     `json:"value"`
}

// Encryption contains the type, used key and the nonce (needed for AES GCM)
type Encryption struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Nonce string `json:"nonce"`
}

// KeyPair represents a public/private key pair
type KeyPair struct {
	key       *rsa.PrivateKey
	encrypter Encrypter
	decrypter Decrypter
}

// NewHybridEncryptionValue returns a new HybridEncryptionValue object for the given parameters
func NewHybridEncryptionValue(encryptionAlgorithm string, encryptedKey string, nonce []byte, encryptedValue []byte) HybridEncryptionValue {
	nonceB64 := base64.StdEncoding.EncodeToString(nonce)
	encryptedValueB64 := base64.StdEncoding.EncodeToString(encryptedValue)

	encryption := &Encryption{Type: encryptionAlgorithm, Key: encryptedKey, Nonce: nonceB64}
	return HybridEncryptionValue{Encryption: *encryption, Value: encryptedValueB64}
}

// Public returns the public key
func (kp *KeyPair) Public() *PublicKey {
	return &PublicKey{&kp.key.PublicKey, kp.encrypter}
}

// Private returns the private key
func (kp *KeyPair) Private() *PrivateKey {
	return &PrivateKey{kp.key, kp.decrypter}
}

// Key interface defines the common functions of a key
type Key interface {
	// AsString returns the key as pem formatted string
	AsString() (string, error)
	// AsBytes returns the key as pem formatted byte array
	AsBytes() ([]byte, error)
	// ToFile writes the key to disk in pem format
	ToFile(path string) error
}
