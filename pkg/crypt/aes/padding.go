package aes

import (
	"bytes"
	"errors"
)

var ErrUnPadding = errors.New("UnPadding error")

const Pkcs5Padding = "PKCS5"
const Pkcs7Padding = "PKCS7"
const ZerosPadding = "ZEROS"

func Padding(padding string, src []byte, blockSize int) []byte {
	switch padding {
	case Pkcs5Padding:
		src = PKCS5Padding(src, blockSize)
	case Pkcs7Padding:
		src = PKCS7Padding(src, blockSize)
	case ZerosPadding:
		src = ZeroPadding(src, blockSize)
	}
	return src
}

func UnPadding(padding string, src []byte) ([]byte, error) {
	switch padding {
	case Pkcs5Padding:
		return PKCS5UnPadding(src)
	case Pkcs7Padding:
		return PKCS7UnPadding(src)
	case ZerosPadding:
		return ZeroUnPadding(src)
	}
	return src, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	return PKCS7Padding(src, blockSize)
}

func PKCS5UnPadding(src []byte) ([]byte, error) {
	return PKCS7UnPadding(src)
}

func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return src, ErrUnPadding
	}
	unPadding := int(src[length-1])
	if length < unPadding {
		return src, ErrUnPadding
	}
	return src[:(length - unPadding)], nil
}

func ZeroPadding(src []byte, blockSize int) []byte {
	paddingCount := blockSize - len(src)%blockSize
	if paddingCount == 0 {
		return src
	} else {
		return append(src, bytes.Repeat([]byte{byte(0)}, paddingCount)...)
	}
}

func ZeroUnPadding(src []byte) ([]byte, error) {
	for i := len(src) - 1; ; i-- {
		if src[i] != 0 {
			return src[:i+1], nil
		}
	}
}
