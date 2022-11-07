package credentials

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/util"
)

const (
	key   = ".credentials.key"
	store = ".credentials.store"
)

var keyPrefix = []byte{
	0x17, 0xe4, 0x27, 0xb7, 0xac, 0xb5, 0x5, 0x7d, 0x37, 0x97, 0x44, 0x8b, 0xe5,
	0xfc, 0xc, 0x6,
}

func newSimpleStore(directory string) (*simpleStore, error) {
	dir := directory
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	key, err := getSecretKey(dir)
	if err != nil {
		return nil, err
	}

	var credentials map[string]*core.Credentials
	storePath := dir + store
	if util.Exists(storePath) {
		credentials, err = readStore(key, storePath)
		if err != nil {
			return nil, err
		}
	} else {
		credentials = make(map[string]*core.Credentials)
	}

	return &simpleStore{
		key:         key,
		store:       storePath,
		credentials: credentials,
	}, nil
}

func getSecretKey(directory string) ([]byte, error) {
	var secretKey []byte
	var err error
	keyPath := directory + key
	if util.Exists(keyPath) {
		secretKey, err = ioutil.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file: %w", err)
		}
	} else {
		secretKey = make([]byte, 16)
		_, err = rand.Read(secretKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create random key: %w", err)
		}

		if !util.Exists(directory) {
			err = os.MkdirAll(directory, 0700)
			if err != nil {
				return nil, fmt.Errorf("failed to create directory for credential store: %w", err)
			}
		}

		err = ioutil.WriteFile(keyPath, secretKey, 0700)
		if err != nil {
			return nil, fmt.Errorf("failed to write random key: %w", err)
		}
	}
	return append(keyPrefix, secretKey...), err
}

func readStore(secretKey []byte, storePath string) (map[string]*core.Credentials, error) {
	ciphertext, err := ioutil.ReadFile(storePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credential store: %w", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher from secret key: %w", err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return readCredentials(ciphertext)
}

func readCredentials(data []byte) (map[string]*core.Credentials, error) {
	credentials := make(map[string]*core.Credentials)
	err := json.Unmarshal(data, &credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret store: %w", err)
	}
	return credentials, nil
}

type simpleStore struct {
	key         []byte
	store       string
	credentials map[string]*core.Credentials
}

// Add adds credentials to the store with the given id.
func (scs *simpleStore) Add(id string, creds *core.Credentials) error {
	scs.credentials[id] = creds
	return scs.writeCredentials()
}

// Remove removes credentials with the given id.
func (scs *simpleStore) Remove(id string) error {
	delete(scs.credentials, id)
	return scs.writeCredentials()
}

// Get returns credentials for the given id.
func (scs *simpleStore) Get(id string) *core.Credentials {
	return scs.credentials[id]
}

func (scs *simpleStore) writeCredentials() error {
	block, err := aes.NewCipher(scs.key)
	if err != nil {
		return fmt.Errorf("failed to create cipher from secret key: %w", err)
	}

	plaintext, err := json.Marshal(&scs.credentials)
	if err != nil {
		return fmt.Errorf("failed to marshall credentials: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// iv =  initialization vector
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return fmt.Errorf("failed to create random iv: %w", err)
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	err = ioutil.WriteFile(scs.store, ciphertext, 0700)
	if err != nil {
		return fmt.Errorf("failed to write credential store: %w", err)
	}
	return nil
}
