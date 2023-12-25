package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// MD5String md5 string
func MD5String(str string) string {
	return MD5Bytes([]byte(str))
}

// MD5Bytes md5 bytes
func MD5Bytes(data []byte) string {
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

// MD5File returns file's hash string
func MD5File(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// MD5YAML 从任意结构体转成yaml后做md5
func MD5YAML(v interface{}) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	m := md5.New()
	m.Write(b)
	return hex.EncodeToString(m.Sum(nil)), nil
}

// MD5JSON 从任意结构体转成json后做md5
func MD5JSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	m := md5.New()
	m.Write(b)
	return hex.EncodeToString(m.Sum(nil)), nil
}
