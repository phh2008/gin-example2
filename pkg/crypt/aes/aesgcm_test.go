package aes

import (
	"testing"

	"com.example/example/pkg/xstring"
)

var idWork32 = func() string { return xstring.RandStr(32) }
var idWork16 = func() string { return xstring.RandStr(16) }

var key string
var iv string

func init() {
	key = idWork32()
	iv = idWork16()
}

func TestAesGcm(t *testing.T) {
	key := []byte(key)
	iv := []byte(iv)
	originPlainText := []byte("abc123")
	aesGcm := NewGCM(key, iv...)
	cipherText, err := aesGcm.EncryptBase64(originPlainText)
	if err != nil {
		t.Error(err)
	}
	t.Log("加密：", cipherText)

	plainText, err := aesGcm.DecryptBase64(cipherText)
	if err != nil || string(plainText) != string(originPlainText) {
		t.Error(err)
	}
	t.Log("解密：", string(plainText))

	cipherText, err = aesGcm.EncryptBase64WithIV(originPlainText)
	if err != nil {
		t.Error(err)
	}
	t.Log("加密WithIV：", cipherText)

	plainText, err = aesGcm.DecryptBase64WithIV(cipherText)
	if err != nil || string(plainText) != string(originPlainText) {
		t.Error(err)
	}
	t.Log("解密WithIV：", string(plainText))
}
