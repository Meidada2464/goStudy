/**
 * Package utils
 * @Author iFurySt <fuqiang.li@baishan.com>
 * @Date 2023/2/6
 */

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

func AESDecryptFromFile(key []byte, filename string) ([]byte, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return AESDecrypt(key, b)
}

func AESDecryptToFile(key, plaintext []byte, filename string) error {
	b, err := AESDecrypt(key, plaintext)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, b, 0644)
}

func AESDecrypt(key, ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func AESEncryptFromFile(key []byte, filename string) ([]byte, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return AESEncrypt(key, b)
}

func AESEncryptToFile(key, plaintext []byte, filename string) error {
	b, err := AESEncrypt(key, plaintext)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, b, 0644)
}

func AESEncrypt(key, plaintext []byte) ([]byte, error) {
	// generate a new aes cipher using 16, 24 or 32 byte long key
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}
