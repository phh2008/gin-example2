package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"com.example/example/pkg/crypt"
)

// GCM aes gcm
type GCM struct {
	Key []byte
	IV  []byte
}

// NewGCM new aes gcm
func NewGCM(key []byte, iv ...byte) *GCM {
	if len(iv) == 0 {
		iv = key[:12]
	}
	return &GCM{
		Key: key,
		IV:  iv,
	}
}

// encrypt 加密
func (obj *GCM) encrypt(plainText []byte, withIV bool) ([]byte, error) {
	if len(obj.Key) != 16 && len(obj.Key) != 24 && len(obj.Key) != 32 {
		return nil, crypt.ErrKeyLength
	}

	block, err := aes.NewCipher(obj.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(obj.IV) != gcm.NonceSize() {
		gcm, err = cipher.NewGCMWithNonceSize(block, len(obj.IV))
		if err != nil {
			return nil, err
		}
	}

	var cipherText []byte
	if withIV {
		cipherText = gcm.Seal(obj.IV, obj.IV, plainText, nil)
	} else {
		cipherText = gcm.Seal(nil, obj.IV, plainText, nil)
	}

	return cipherText, nil
}

// EncryptWithIV 加密,带有iv前缀
func (obj *GCM) EncryptWithIV(plainText []byte) ([]byte, error) {
	return obj.encrypt(plainText, true)
}

// Encrypt 加密，无前缀
func (obj *GCM) Encrypt(plainText []byte) ([]byte, error) {
	return obj.encrypt(plainText, false)
}

// EncryptBase64WithIV 加密base64，带iv前缀
func (obj *GCM) EncryptBase64WithIV(plainText []byte) (string, error) {
	cipherText, err := obj.EncryptWithIV(plainText)
	if err != nil {
		return "", err
	}
	ret := base64.StdEncoding.EncodeToString(cipherText)
	return ret, nil
}

// EncryptBase64 加密base64
func (obj *GCM) EncryptBase64(plainText []byte) (string, error) {
	cipherText, err := obj.Encrypt(plainText)
	if err != nil {
		return "", err
	}
	ret := base64.StdEncoding.EncodeToString(cipherText)
	return ret, nil
}

// decrypt 解密
func (obj *GCM) decrypt(cipherText []byte, withIV bool) ([]byte, error) {
	if len(obj.Key) != 16 && len(obj.Key) != 24 && len(obj.Key) != 32 {
		return nil, crypt.ErrKeyLength
	}
	if len(cipherText) < aes.BlockSize {
		return nil, crypt.ErrCipherText
	}

	block, err := aes.NewCipher(obj.Key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(obj.IV) != gcm.NonceSize() {
		gcm, err = cipher.NewGCMWithNonceSize(block, len(obj.IV))
		if err != nil {
			return nil, err
		}
	}

	var plainText []byte
	if withIV {
		iv, cipherText := cipherText[:gcm.NonceSize()], cipherText[gcm.NonceSize():]
		plainText, err = gcm.Open(nil, iv, cipherText, nil)
	} else {
		plainText, err = gcm.Open(nil, obj.IV, cipherText, nil)
	}
	if err != nil {
		return nil, err
	}

	return plainText, err
}

// DecryptWithIV 解密，带有iv前缀
func (obj *GCM) DecryptWithIV(cipherText []byte) ([]byte, error) {
	return obj.decrypt(cipherText, true)
}

// Decrypt 解密
func (obj *GCM) Decrypt(cipherText []byte) ([]byte, error) {
	return obj.decrypt(cipherText, false)
}

// DecryptBase64WithIV 解密base64，带有iv前缀
func (obj *GCM) DecryptBase64WithIV(cipherStr string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherStr)
	if err != nil {
		return nil, err
	}
	return obj.DecryptWithIV(cipherText)
}

// DecryptBase64 解密base64
func (obj *GCM) DecryptBase64(cipherStr string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherStr)
	if err != nil {
		return nil, err
	}
	return obj.Decrypt(cipherText)
}
