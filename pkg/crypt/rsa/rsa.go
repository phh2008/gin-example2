package rsa

import (
	"crypto"
)

type Key interface {
	~string | ~[]byte
}

// Encrypt encrypts by rsa with public key or private key.
// 通过 rsa 公钥或私钥加密
func Encrypt[T Key](rsaKey T, src []byte) CodecResult {
	keyPair := NewKeyPair()
	keyPair.SetPublicKey([]byte(rsaKey))
	keyPair.SetPrivateKey([]byte(rsaKey))
	var c CodecResult
	// 私钥加密
	if keyPair.IsPrivateKey() {
		c.Raw, c.Error = keyPair.EncryptByPrivateKey(src)
		return c
	}
	// 公钥加密
	c.Raw, c.Error = keyPair.EncryptByPublicKey(src)
	return c
}

// Decrypt decrypts by rsa with private key or public key.
// 通过 rsa 私钥或公钥解密
func Decrypt[T Key](rsaKey T, src []byte) CodecResult {
	keyPair := NewKeyPair()
	keyPair.SetPublicKey([]byte(rsaKey))
	keyPair.SetPrivateKey([]byte(rsaKey))
	var p CodecResult
	if keyPair.IsPublicKey() {
		p.Raw, p.Error = keyPair.DecryptByPublicKey(src)
		return p
	}
	p.Raw, p.Error = keyPair.DecryptByPrivateKey(src)
	return p
}

// CreateSign by rsa with private key.
// 通过 rsa 私钥签名
func CreateSign[T Key](privateKey T, hash crypto.Hash, src []byte) CodecResult {
	keyPair := NewKeyPair()
	keyPair.SetPrivateKey([]byte(privateKey))
	keyPair.SetHash(hash)
	var s CodecResult
	s.Raw, s.Error = keyPair.SignByPrivateKey(src)
	return s
}

// VerifySign verify sign by rsa with public key.
// 通过 rsa 公钥验签
func VerifySign[T Key](publicKey T, hash crypto.Hash, src []byte, sign []byte) (bool, error) {
	if len(src) == 0 || len(sign) == 0 {
		return false, nil
	}
	keyPair := NewKeyPair()
	keyPair.SetPublicKey([]byte(publicKey))
	keyPair.SetHash(hash)
	err := keyPair.VerifyByPublicKey(src, sign)
	return err == nil, err
}
