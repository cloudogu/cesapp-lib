package keys

import (
	"crypto/rsa"
	"crypto/sha256"
	"io"
)

func init() {
	oaespProvider := newOaespKeyProvider()
	register("oaesp", oaespProvider)
	register("", oaespProvider)
}

func newOaespKeyProvider() *KeyProvider {
	return &KeyProvider{oaespEncrypter, oaespDecrypter}
}

const label string = "ces-config"

func oaespEncrypter(random io.Reader, pub *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), random, pub, msg, []byte(label))
}

func oaespDecrypter(random io.Reader, priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), random, priv, ciphertext, []byte(label))
}
