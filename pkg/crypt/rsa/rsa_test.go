package rsa

import (
	"crypto"
	"testing"

	"com.example/example/pkg/crypt/openssl"
	"github.com/stretchr/testify/assert"
)

func TestRsa_EncryptWithPublicKey(t *testing.T) {
	input := []byte("hello world!!!")
	pubKey, priKey, _ := openssl.RSA.GenKeyPair(openssl.PKCS8, 2048)
	t.Log("公钥：\n", string(pubKey))
	t.Log("私钥：\n", string(priKey))

	enMsg := Encrypt(pubKey, input)
	assert.Nil(t, enMsg.Error)
	t.Log("加密后：", enMsg.ToBase64String())

	deMsg := Decrypt(priKey, enMsg.Raw)
	assert.Nil(t, deMsg.Error)
	t.Log("解密后：", deMsg.String())
}

func TestRsa_EncryptWithPrivateKey(t *testing.T) {
	input := []byte("hello world!!!")
	pubKey, priKey, _ := openssl.RSA.GenKeyPair(openssl.PKCS8, 2048)
	t.Log("公钥：\n", string(pubKey))
	t.Log("私钥：\n", string(priKey))

	enMsg := Encrypt(priKey, input)
	assert.Nil(t, enMsg.Error)
	t.Log("加密后：", enMsg.ToBase64String())

	deMsg := Decrypt(pubKey, enMsg.Raw)
	assert.Nil(t, deMsg.Error)
	t.Log("解密后：", deMsg.String())
}

func TestRsa_SignWithPrivateKy(t *testing.T) {
	input := []byte("hello world!!!")
	pubKey, priKey, _ := openssl.RSA.GenKeyPair(openssl.PKCS8, 2048)
	t.Log("公钥：\n", string(pubKey))
	t.Log("私钥：\n", string(priKey))

	signResult := CreateSign(priKey, crypto.SHA256, input)
	t.Log("签名串：", signResult.ToBase64String())
	assert.Nil(t, signResult.Error)
	ok, err := VerifySign(pubKey, crypto.SHA256, input, signResult.Raw)
	assert.Nil(t, err)
	t.Log("签名验证是否通过：", ok)
}
