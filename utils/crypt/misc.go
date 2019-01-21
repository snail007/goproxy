package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

type CryptTool struct{}

var CryptTools = NewCryptTool()

func NewCryptTool() *CryptTool {
	return &CryptTool{}
}

func (encrypt *CryptTool) Base64Encode(str string) string {
	return string([]byte(base64.StdEncoding.EncodeToString([]byte(str))))
}

func (encrypt *CryptTool) Base64EncodeBytes(bytes []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(bytes))
}

func (encrypt *CryptTool) Base64Decode(str string) (string, error) {
	by, err := base64.StdEncoding.DecodeString(str)
	return string(by), err
}

func (encrypt *CryptTool) Base64DecodeBytes(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

func (encrypt *CryptTool) MD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
