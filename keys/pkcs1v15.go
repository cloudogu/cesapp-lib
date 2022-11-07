package keys

import (
	"crypto/rsa"
	"io"
)

func init() {
	pkcs1v15Provider := newPkcs1v15KeyProvider()
	register("pkcs1v15", pkcs1v15Provider)
}

func newPkcs1v15KeyProvider() *KeyProvider {
	return &KeyProvider{pkcs1v15Encrypter, pkcs1v15Decrypter}
}

func pkcs1v15Encrypter(random io.Reader, pub *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(random, pub, msg)
}

func pkcs1v15Decrypter(random io.Reader, priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(random, priv, ciphertext)
}
