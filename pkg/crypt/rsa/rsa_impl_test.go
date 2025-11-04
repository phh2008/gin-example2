package rsa

import (
	"crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	rsaInput        = "hello world"
	pkcs1PublicKey  = ``
	pkcs1PrivateKey = ``
	pkcs8PublicKey  = ``
	pkcs8PrivateKey = ``
)

func TestRsa_PKCS1_Encrypt(t *testing.T) {
	keyPair := NewKeyPair()
	keyPair.SetPublicKey([]byte(pkcs1PublicKey))
	keyPair.SetPrivateKey([]byte(pkcs1PrivateKey))

	dst1, err1 := keyPair.EncryptByPublicKey([]byte(rsaInput))
	assert.Nil(t, err1)
	dst2, err2 := keyPair.DecryptByPrivateKey(dst1)
	assert.Nil(t, err2)
	assert.Equal(t, []byte(rsaInput), dst2)

	dst3, err3 := keyPair.EncryptByPrivateKey([]byte(rsaInput))
	assert.Nil(t, err3)
	dst4, err4 := keyPair.DecryptByPublicKey(dst3)
	assert.Nil(t, err4)
	assert.Equal(t, []byte(rsaInput), dst4)
}

func TestRsa_PKCS1_Sign(t *testing.T) {
	keyPair := NewKeyPair()
	keyPair.SetPublicKey([]byte(pkcs1PublicKey))
	keyPair.SetPrivateKey([]byte(pkcs1PrivateKey))
	keyPair.SetHash(crypto.SHA224)

	dst1, err1 := keyPair.SignByPrivateKey([]byte(rsaInput))
	assert.Nil(t, err1)
	err2 := keyPair.VerifyByPublicKey([]byte(rsaInput), dst1)
	assert.Nil(t, err2)
	assert.Equal(t, nil, err2)
}

func TestRsa_PKCS8_Encrypt(t *testing.T) {
	keyPair := NewKeyPair()
	keyPair.SetPublicKey([]byte(pkcs8PublicKey))
	keyPair.SetPrivateKey([]byte(pkcs8PrivateKey))

	dst1, err1 := keyPair.EncryptByPublicKey([]byte(rsaInput))
	assert.Nil(t, err1)
	dst2, err2 := keyPair.DecryptByPrivateKey(dst1)
	assert.Nil(t, err2)
	assert.Equal(t, []byte(rsaInput), dst2)

	dst3, err3 := keyPair.EncryptByPrivateKey([]byte(rsaInput))
	assert.Nil(t, err3)
	dst4, err4 := keyPair.DecryptByPublicKey(dst3)
	assert.Nil(t, err4)
	assert.Equal(t, []byte(rsaInput), dst4)
}

func TestRsa_PKCS8_Sign(t *testing.T) {
	keyPair := NewKeyPair()
	keyPair.SetPublicKey([]byte(pkcs8PublicKey))
	keyPair.SetPrivateKey([]byte(pkcs8PrivateKey))
	keyPair.SetHash(crypto.SHA384)

	dst1, err1 := keyPair.SignByPrivateKey([]byte(rsaInput))
	assert.Nil(t, err1)
	err2 := keyPair.VerifyByPublicKey([]byte(rsaInput), dst1)
	assert.Nil(t, err2)
	assert.Equal(t, nil, err2)
}

func TestRsa_IsKey(t *testing.T) {
	keyPair := NewKeyPair()

	keyPair.SetPublicKey([]byte(pkcs1PublicKey))
	keyPair.SetPrivateKey([]byte(pkcs1PrivateKey))
	assert.Equal(t, true, keyPair.IsPublicKey())
	assert.Equal(t, true, keyPair.IsPrivateKey())

	keyPair.SetPublicKey([]byte(pkcs8PublicKey))
	keyPair.SetPrivateKey([]byte(pkcs8PrivateKey))
	assert.Equal(t, true, keyPair.IsPublicKey())
	assert.Equal(t, true, keyPair.IsPrivateKey())

	keyPair.SetPublicKey([]byte(pkcs1PrivateKey))
	keyPair.SetPrivateKey([]byte(pkcs1PublicKey))
	assert.Equal(t, false, keyPair.IsPublicKey())
	assert.Equal(t, false, keyPair.IsPrivateKey())

	keyPair.SetPublicKey([]byte(pkcs8PrivateKey))
	keyPair.SetPrivateKey([]byte(pkcs8PublicKey))
	assert.Equal(t, false, keyPair.IsPublicKey())
	assert.Equal(t, false, keyPair.IsPrivateKey())

	keyPair.SetPublicKey([]byte(""))
	keyPair.SetPrivateKey([]byte(""))
	assert.Equal(t, false, keyPair.IsPublicKey())
	assert.Equal(t, false, keyPair.IsPrivateKey())

	keyPair.SetPublicKey([]byte("xxx"))
	keyPair.SetPrivateKey([]byte("xxx"))
	assert.Equal(t, false, keyPair.IsPublicKey())
	assert.Equal(t, false, keyPair.IsPrivateKey())
}

func TestRsa_Empty_Src(t *testing.T) {
	keyPair := NewKeyPair()
	keyPair.SetPublicKey([]byte(pkcs1PublicKey))
	keyPair.SetPrivateKey([]byte(pkcs1PrivateKey))

	empty := ""
	dst1, err1 := keyPair.EncryptByPublicKey([]byte(empty))
	assert.Nil(t, err1)
	dst2, err2 := keyPair.DecryptByPrivateKey(dst1)
	assert.Nil(t, err2)
	assert.Equal(t, []byte(empty), dst2)

	dst3, err3 := keyPair.EncryptByPrivateKey([]byte(empty))
	assert.Nil(t, err3)
	dst4, err4 := keyPair.DecryptByPublicKey(dst3)
	assert.Nil(t, err4)
	assert.Equal(t, []byte(empty), dst4)
}

func TestRsa_PublicKey_Error(t *testing.T) {
	invalidPublicKey := `-----BEGIN PUBLIC KEY-----
xxxx
-----END PUBLIC KEY-----`

	keyPair := NewKeyPair()
	keyPair.SetPrivateKey([]byte(pkcs1PrivateKey))
	keyPair.SetHash(crypto.SHA1)

	sign, _ := keyPair.SignByPrivateKey([]byte(rsaInput))

	keyPair.SetPublicKey([]byte(pkcs1PrivateKey))
	_, err1 := keyPair.EncryptByPublicKey([]byte(rsaInput))
	assert.Equal(t, invalidPublicKeyError(), err1)
	_, err2 := keyPair.DecryptByPublicKey(sign)
	assert.Equal(t, invalidPublicKeyError(), err2)
	err3 := keyPair.VerifyByPublicKey([]byte(rsaInput), sign)
	assert.Equal(t, invalidPublicKeyError(), err3)

	keyPair.SetPublicKey([]byte(invalidPublicKey))
	_, err4 := keyPair.EncryptByPublicKey([]byte(rsaInput))
	assert.Equal(t, invalidPublicKeyError(), err4)
	_, err5 := keyPair.DecryptByPublicKey(sign)
	assert.Equal(t, invalidPublicKeyError(), err5)
	err6 := keyPair.VerifyByPublicKey([]byte(rsaInput), sign)
	assert.Equal(t, invalidPublicKeyError(), err6)

	keyPair.SetPrivateKey([]byte(pkcs1PrivateKey))
	_, err7 := keyPair.EncryptByPublicKey([]byte(rsaInput))
	assert.Equal(t, invalidPublicKeyError(), err7)
	_, err8 := keyPair.DecryptByPublicKey([]byte(rsaInput))
	assert.Equal(t, invalidPublicKeyError(), err8)
	err9 := keyPair.VerifyByPublicKey([]byte(rsaInput), sign)
	assert.Equal(t, invalidPublicKeyError(), err9)
}

func TestRsa_PrivateKey_Error(t *testing.T) {
	invalidPrivateKey := `-----BEGIN PRIVATE KEY-----
xxxx
-----END PRIVATE KEY-----`

	keyPair := NewKeyPair()
	keyPair.SetHash(crypto.SHA512)

	keyPair.SetPrivateKey([]byte(invalidPrivateKey))
	_, err1 := keyPair.EncryptByPrivateKey([]byte(rsaInput))
	assert.Equal(t, invalidPrivateKeyError(), err1)
	_, err2 := keyPair.DecryptByPrivateKey([]byte(rsaInput))
	assert.Equal(t, invalidPrivateKeyError(), err2)
	_, err3 := keyPair.SignByPrivateKey([]byte(rsaInput))
	assert.Equal(t, invalidPrivateKeyError(), err3)

	keyPair.SetPrivateKey([]byte(pkcs1PublicKey))
	_, err4 := keyPair.EncryptByPrivateKey([]byte(rsaInput))
	assert.Equal(t, invalidPrivateKeyError(), err4)
	_, err5 := keyPair.DecryptByPrivateKey([]byte(rsaInput))
	assert.Equal(t, invalidPrivateKeyError(), err5)
	_, err6 := keyPair.SignByPrivateKey([]byte(rsaInput))
	assert.Equal(t, invalidPrivateKeyError(), err6)

	keyPair.SetPrivateKey([]byte(pkcs1PublicKey))
	_, err7 := keyPair.EncryptByPrivateKey([]byte(rsaInput))
	assert.Equal(t, invalidPrivateKeyError(), err7)
	_, err8 := keyPair.DecryptByPrivateKey([]byte(rsaInput))
	assert.Equal(t, invalidPrivateKeyError(), err8)
	_, err9 := keyPair.SignByPrivateKey([]byte(rsaInput))
	assert.Equal(t, invalidPrivateKeyError(), err9)
}

func TestRsa_Hash_Error(t *testing.T) {
	keyPair := NewKeyPair()
	keyPair.SetPublicKey([]byte(pkcs1PublicKey))
	keyPair.SetPrivateKey([]byte(pkcs1PrivateKey))
	keyPair.SetHash(crypto.MD4)

	_, err1 := keyPair.SignByPrivateKey([]byte(rsaInput))
	assert.Equal(t, unsupportedHashError(), err1)
	err2 := keyPair.VerifyByPublicKey([]byte(rsaInput), []byte(rsaInput))
	assert.Equal(t, unsupportedHashError(), err2)
}
