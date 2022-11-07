package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"

	"crypto/cipher"

	"crypto/aes"

	"encoding/json"

	"fmt"

	"github.com/pkg/errors"
)

// PublicKey is the public key part of the KeyPair.
type PublicKey struct {
	key       *rsa.PublicKey
	encrypter Encrypter
}

// AsBytes returns the key as pem formatted byte array.
func (pk *PublicKey) AsBytes() ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(pk.key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key")
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}

	return pem.EncodeToMemory(pemkey), nil
}

// AsString returns the key as pem formatted string.
func (pk *PublicKey) AsString() (string, error) {
	bytes, err := pk.AsBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

// ToFile writes the key to disk in pem format.
func (pk *PublicKey) ToFile(path string) error {
	bytes, err := pk.AsBytes()
	if err != nil {
		return err
	}

	// TODO check permissions
	err = ioutil.WriteFile(path, bytes, 0744)
	if err != nil {
		return fmt.Errorf("failed to write public key to file %s: %w", path, err)
	}
	return nil
}

// Encrypt encrypts the given input.
// In cases where the input can not be encrypted with RSA because it is too long,
// we switch to a hybrid encryption (i. e. using symmetric crypto for the actual content via a randomly generated key which in turn is encrypted with RSA).
func (pk *PublicKey) Encrypt(input string) (string, error) {
	inputBytes := []byte(input)
	if pk.canEncryptWithRSA(input) {
		return pk.encryptWithRSAandEncode(inputBytes)
	}
	return pk.encryptHybrid(inputBytes)
}

// We chose the length of 64 bytes to ensure the encryption works for different
// RSA padding schemes.
func (pk *PublicKey) canEncryptWithRSA(input string) bool {
	return len(input) <= MaxRSAEncryptionLength
}

func (pk *PublicKey) encryptWithRSAandEncode(input []byte) (string, error) {
	ciphertext, err := pk.encrypter(rand.Reader, pk.key, input)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt input text: %w", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (pk *PublicKey) encryptHybrid(input []byte) (string, error) {
	// Generate Key + Nonce
	aesKey, err := generateCryptographicallySecureRandomBits(AesKeyBitLength)
	if err != nil {
		return "", err
	}
	// We do not ensure the generated nonce is really a number only used once (nonce) here,
	// but consider it sufficiently unlikely to generate the same pair of key + nonce twice.
	nonce, err := generateCryptographicallySecureRandomBits(NonceBitLength)
	if err != nil {
		return "", err
	}

	// symmetrically encrypt the actual value...
	encryptedValue, err := encryptWithAesGcm(aesKey, nonce, input)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt value with AES GCM: %w", err)
	}

	// ... then encrypt the symmetric key regularly (i. e. asymmetrically).
	encryptedAesKey, err := pk.encryptWithRSAandEncode(aesKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt aes key: %w", err)
	}

	// finally,  join both parts for persisting them.
	value := NewHybridEncryptionValue(AesGcm, encryptedAesKey, nonce, encryptedValue)
	valueAsJSON, _ := json.Marshal(value)

	return string(valueAsJSON), nil
}

// Provides random data produced by an underlying Cryptographically Secure Pseudo Random Number Generator.
// Length have to be a multiple of 8 because the random data will be written in a byte array.
func generateCryptographicallySecureRandomBits(length int) ([]byte, error) {
	if length%8 != 0 {
		return nil, fmt.Errorf("length %d is not a multiple of 8", length)
	}
	// divide by 8 to make up for byte vs bits as given by method contract
	randomBytes := make([]byte, length/8)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, errors.Wrapf(err, "failed to generate random")
	}
	return randomBytes, nil
}

// encrypt with AES in GCM mode with the given length of the key.
func encryptWithAesGcm(aesKey []byte, nonce []byte, input []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create aes cipher: %w", err)
	}

	aesGcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm cipher: %w", err)
	}

	return aesGcm.Seal(nil, nonce, input, nil), nil
}
