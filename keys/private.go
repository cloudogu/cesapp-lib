package keys

// PrivateKey is the private key part of the KeyPair
import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"

	"strings"

	"encoding/json"

	"crypto/aes"
	"crypto/cipher"

	"github.com/pkg/errors"
)

type PrivateKey struct {
	key       *rsa.PrivateKey
	decrypter Decrypter
}

// AsBytes returns the key as pem formatted byte array
func (pk *PrivateKey) AsBytes() ([]byte, error) {
	var pemkey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pk.key),
	}

	return pem.EncodeToMemory(pemkey), nil
}

// AsString returns the key as pem formatted string
func (pk *PrivateKey) AsString() (string, error) {
	bytes, err := pk.AsBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

// ToFile writes the key to disk in pem format
func (pk *PrivateKey) ToFile(path string) error {
	bytes, err := pk.AsBytes()
	if err != nil {
		return err
	}

	// TODO check permissions
	err = ioutil.WriteFile(path, bytes, 0744)
	if err != nil {
		return errors.Wrap(err, "failed to write private key to file "+path)
	}
	return nil
}

// Decrypt decrypts a text which was encrypted with the Encrypt function of the
// Public key of the same key pair.
// In cases where the input is a meta value, we have to decrypt the symmetric key and use it for decrypting
// the real value.
func (pk *PrivateKey) Decrypt(input string) (string, error) {
	if isMetaValue(input) {
		return pk.decryptHybrid(input)
	}

	return pk.decryptRSAEncryptedB64String(input)
}

func isMetaValue(input string) bool {
	return strings.HasPrefix(input, "{") && strings.HasSuffix(input, "}")
}

func (pk *PrivateKey) decryptRSAEncryptedB64String(input string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64 input")
	}

	plaintext, err := pk.decrypter(rand.Reader, pk.key, ciphertext)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt input")
	}

	return string(plaintext), nil
}

func (pk *PrivateKey) decryptHybrid(input string) (string, error) {
	hybridEncryptionValue := &HybridEncryptionValue{}

	if err := json.Unmarshal([]byte(input), hybridEncryptionValue); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal HybridEncryptionValue input")
	}

	if hybridEncryptionValue.Encryption.Type == AesGcm {
		return pk.decryptHybridWithAesGcm(*hybridEncryptionValue)
	}
	return "", errors.Errorf("unsupported encryption algorithm detected: %s", hybridEncryptionValue.Encryption.Type)
}

func (pk *PrivateKey) decryptHybridWithAesGcm(hybridEncryptionValue HybridEncryptionValue) (string, error) {

	plainSymmetricKey, err := pk.decryptRSAEncryptedB64String(hybridEncryptionValue.Encryption.Key)
	if err != nil {
		return "", errors.Wrap(err, "could not decrypt hybrid encryption value")
	}

	aesCipher, err := aes.NewCipher([]byte(plainSymmetricKey))
	if err != nil {
		return "", errors.Wrap(err, "failed to create aes cipher")
	}
	aesGcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", errors.Wrap(err, "failed to create aes gcm")
	}
	providedNonce, err := base64.StdEncoding.DecodeString(hybridEncryptionValue.Encryption.Nonce)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode(base64) nonce")
	}
	encryptedValue, err := base64.StdEncoding.DecodeString(hybridEncryptionValue.Value)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode(base64) encrypted value")
	}
	decryptedValue, err := aesGcm.Open(nil, providedNonce, encryptedValue, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to aes-decrypt encryptedValue")
	}

	return string(decryptedValue), nil
}
