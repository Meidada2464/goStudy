package AEScode

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"testing"
)

func TestUncode(t *testing.T) {
	key := []byte{55, 57, 48, 101, 102, 57, 51, 48, 99, 50, 55, 56, 48, 97, 57, 49, 48, 101, 51, 54, 52, 102, 55, 55, 57, 50, 102, 52, 57, 102, 54, 49}
	file, err := os.Open("/Users/meifengfeng/workSpace/study/goStudy/goStudy/AEScode/config.yaml")
	if err != nil {
		return
	}
	buf := make([]byte, 1024*5)
	n, err := file.Read(buf)
	if err != nil {
		fmt.Println("bbb")
		return
	}
	if n != 0 {
		fmt.Println("aaa")
		return
	}

	c, err := AESDecrypt(key, buf)
	if err != nil {
		return
	}
	str := string(c)

	fmt.Println(str)
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
