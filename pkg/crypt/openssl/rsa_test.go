package openssl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	publicKey1  []byte
	privateKey1 []byte
	publicKey8  []byte
	privateKey8 []byte
)

func TestRsa_GenKeyPair(t *testing.T) {
	publicKey1, privateKey1, _ = RSA.GenKeyPair(PKCS1, 1024)
	fmt.Printf("生成 1024 字节 PKCS#1 格式公钥：\n%s", publicKey1)
	fmt.Printf("生成 1024 字节 PKCS#1 格式私钥：\n%s", privateKey1)

	publicKey8, privateKey8, _ = RSA.GenKeyPair(PKCS8, 2048)
	fmt.Printf("生成 2048 字节 PKCS#8 格式公钥：\n%s", publicKey8)
	fmt.Printf("生成 2048 字节 PKCS#8 格式私钥：\n%s", privateKey8)
}

func TestRsa_VerifyKeyPair(t *testing.T) {
	publicKey1, privateKey1, _ = RSA.GenKeyPair(PKCS1, 1024)
	publicKey8, privateKey8, _ = RSA.GenKeyPair(PKCS8, 2048)

	assert.Equal(t, true, RSA.VerifyKeyPair(publicKey1, privateKey1))
	assert.Equal(t, true, RSA.VerifyKeyPair(publicKey8, privateKey8))
	assert.Equal(t, false, RSA.VerifyKeyPair(publicKey1, privateKey8))
	assert.Equal(t, false, RSA.VerifyKeyPair(publicKey8, privateKey1))
	assert.Equal(t, false, RSA.VerifyKeyPair(publicKey1, []byte("xxx")))
	assert.Equal(t, false, RSA.VerifyKeyPair([]byte("xxx"), privateKey1))
}

func TestRsa_IsPublicKey(t *testing.T) {
	publicKey1, privateKey1, _ = RSA.GenKeyPair(PKCS1, 1024)
	publicKey8, privateKey8, _ = RSA.GenKeyPair(PKCS8, 2048)

	assert.Equal(t, true, RSA.IsPublicKey(publicKey1))
	assert.Equal(t, true, RSA.IsPrivateKey(privateKey1))
	assert.Equal(t, false, RSA.IsPublicKey(privateKey1))
	assert.Equal(t, false, RSA.IsPrivateKey(publicKey1))

	assert.Equal(t, true, RSA.IsPublicKey(publicKey8))
	assert.Equal(t, true, RSA.IsPrivateKey(privateKey8))
	assert.Equal(t, false, RSA.IsPublicKey(privateKey8))
	assert.Equal(t, false, RSA.IsPrivateKey(publicKey8))

	assert.Equal(t, false, RSA.IsPublicKey([]byte("xxx")))
	assert.Equal(t, false, RSA.IsPrivateKey([]byte("xxx")))
}

func TestRsa_ParsePublicKey(t *testing.T) {
	publicKey1, privateKey1, _ = RSA.GenKeyPair(PKCS1, 1024)
	publicKey8, privateKey8, _ = RSA.GenKeyPair(PKCS8, 2048)

	pub1, err1 := RSA.ParsePublicKey(publicKey1)
	assert.Nil(t, err1)
	pri1, err1 := RSA.ParsePrivateKey(privateKey1)
	assert.Nil(t, err1)
	assert.Equal(t, pub1, &pri1.PublicKey)

	pub8, err8 := RSA.ParsePublicKey(publicKey8)
	assert.Nil(t, err8)
	pri8, err8 := RSA.ParsePrivateKey(privateKey8)
	assert.Nil(t, err8)
	assert.Equal(t, pub8, &pri8.PublicKey)
}

func TestRsa_ExportPublicKey(t *testing.T) {
	publicKey1, privateKey1, _ = RSA.GenKeyPair(PKCS1, 1024)
	publicKey8, privateKey8, _ = RSA.GenKeyPair(PKCS8, 2048)

	actual1, err1 := RSA.ExportPublicKey(privateKey1)
	assert.Nil(t, err1)
	assert.Equal(t, publicKey1, actual1)

	actual8, err8 := RSA.ExportPublicKey(privateKey8)
	assert.Nil(t, err8)
	assert.Equal(t, publicKey8, actual8)
}

func TestRsaFmtKey(t *testing.T) {
	publicKey1, privateKey1, _ = RSA.GenKeyPair(PKCS1, 1024)
	pubKey, err := RSA.CompressKey(publicKey1)
	assert.Nil(t, err)
	t.Log("压缩后的pubKey:\n", string(pubKey))
	priKey, err := RSA.CompressKey(privateKey1)
	assert.Nil(t, err)
	t.Log("压缩后的priKey:\n", string(priKey))

	priKey = RSA.FormatPrivateKey(PKCS1, priKey)
	t.Log("格式化后的priKey:\n", string(priKey))

	pubKey = RSA.FormatPublicKey(PKCS1, pubKey)
	t.Log("格式化后的priKey:\n", string(pubKey))
}

func TestRsaFmtKey2(t *testing.T) {
	publicKey1, privateKey1, _ = RSA.GenKeyPair(PKCS8, 2048)
	pubKey, err := RSA.CompressKey(publicKey1)
	assert.Nil(t, err)
	t.Log("压缩后的pubKey:\n", string(pubKey))
	priKey, err := RSA.CompressKey(privateKey1)
	assert.Nil(t, err)
	t.Log("压缩后的priKey:\n", string(priKey))

	priKey = RSA.FormatPrivateKey(PKCS8, priKey)
	t.Log("格式化后的priKey:\n", string(priKey))

	pubKey = RSA.FormatPublicKey(PKCS8, pubKey)
	t.Log("格式化后的pubKey:\n", string(pubKey))
}
