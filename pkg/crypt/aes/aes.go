package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func ECBEncrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	return ecbEncrypt(block, src, padding)
}

func ECBDecrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	return ecbDecrypt(block, src, padding)
}

func CBCEncrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cbcEncrypt(block, src, iv, padding)
}

func CBCDecrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cbcDecrypt(block, src, iv, padding)
}

func NewCipher(key []byte) (cipher.Block, error) {
	return aes.NewCipher(aesKeyPending(key))
}

// aesKeyPending The length of the key can be 16/24/32 characters (128/192/256 bits)
func aesKeyPending(key []byte) []byte {
	k := len(key)
	count := 0
	switch true {
	case k <= 16:
		count = 16 - k
	case k <= 24:
		count = 24 - k
	case k <= 32:
		count = 32 - k
	default:
		return key[:32]
	}
	if count == 0 {
		return key
	}
	return append(key, bytes.Repeat([]byte{0}, count)...)
}
