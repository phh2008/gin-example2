package aes

import (
	"encoding/base64"
	"testing"
)

func TestEncrypt(t *testing.T) {
	var key = []byte(idWork32())
	var iv = []byte(idWork16())
	t.Log("key: ", string(key))
	t.Log("iv: ", string(iv))

	var src = "abc123_!@#"
	enc, err := CBCEncrypt([]byte(src), key, iv, Pkcs5Padding)
	if err != nil {
		panic(err)
	}
	t.Log("CBC加密：", base64.StdEncoding.EncodeToString(enc))
	str, err := CBCDecrypt(enc, key, iv, Pkcs5Padding)
	if err != nil {
		panic(err)
	}
	t.Log("CBC解密：", string(str))
}

func TestECBEncrypt(t *testing.T) {
	// AES/ECB/PKCS5Padding
	var key = []byte(idWork32())
	t.Log("key: ", string(key))

	var src = []byte("abc123")
	enc, err := ECBEncrypt(src, key, Pkcs5Padding)
	if err != nil {
		panic(err)
	}
	t.Log("ECB加密：", base64.StdEncoding.EncodeToString(enc))
	str, err := ECBDecrypt(enc, key, Pkcs5Padding)
	if err != nil {
		panic(err)
	}
	t.Log("ECB解密：", string(str))
}
