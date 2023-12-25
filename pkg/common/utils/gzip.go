package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
)

// Gzip compress bytes to bytes
func Gzip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, len(data)/3))
	gz := gzip.NewWriter(buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GzipJSON encodes value as compressed json bytes to reader
func GzipJSON(value interface{}) (io.Reader, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}

// GzipJSONBytes encodes values as compressed json and returns bytes
func GzipJSONBytes(value interface{}) ([]byte, error) {
	reader, err := GzipJSON(value)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(reader)
}

// UnGzipJSON decodes compressed json bytes in reader to object
func UnGzipJSON(reader io.Reader, value interface{}) error {
	rd, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer rd.Close()

	decoder := json.NewDecoder(rd)
	return decoder.Decode(value)
}

// UnGzipJSONBytes decodes compressd json bytes to object
func UnGzipJSONBytes(data []byte, value interface{}) error {
	data, err := UnGzip(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// UnGzip 解gzip数据
func UnGzip(data []byte) ([]byte, error) {
	rd, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	return io.ReadAll(rd)
}

// JSONToFile 按照 json 格式写入数据到文件
func JSONToFile(file string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(file, b, 0644)
}

// JSONFromFile 按照 json 格式读入数据
func JSONFromFile(file string, data interface{}) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, data)
}
