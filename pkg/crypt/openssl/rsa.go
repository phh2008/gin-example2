package openssl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
)

var (
	// returns an invalid RSA public key error.
	// 返回无效的 RSA 公钥错误
	invalidRsaPublicKeyError = func() error {
		return fmt.Errorf("openssl/rsa: invalid rsa public key, please make sure the public key is valid")
	}
	// returns an invalid RSA private key error.
	// 返回无效的 RSA 私钥错误
	invalidRSAPrivateKeyError = func() error {
		return fmt.Errorf("openssl/rsa: invalid rsa private key, please make sure the private key is valid")
	}
	// returns an invalid RSA key error.
	// 返回无效的 RSA 密钥错误
	invalidRSAKeyError = func() error {
		return fmt.Errorf("openssl/rsa: invalid rsa key, please make sure the key is valid")
	}
)

// rsaKeyPair defines a rsaKeyPair struct.
// 定义 rsaKeyPair 结构体
type rsaKeyPair struct {
	publicKey  []byte
	privateKey []byte
}

// newRsaKeyPair returns a new rsaKeyPair instance.
// 初始化 rsaKeyPair 结构体
func newRsaKeyPair() *rsaKeyPair {
	return &rsaKeyPair{}
}

// GenKeyPair generates key pair.
// 生成密钥对
func (r rsaKeyPair) GenKeyPair(pkcs keyFormat, bits int) (publicKey, privateKey []byte, err error) {
	var pri *rsa.PrivateKey
	pri, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	if pkcs == PKCS1 {
		privateBytes := x509.MarshalPKCS1PrivateKey(pri)
		privateKey = pem.EncodeToMemory(&pem.Block{
			Type:  BlockTypeRsaPrivateKey,
			Bytes: privateBytes,
		})
		publicBytes := x509.MarshalPKCS1PublicKey(&pri.PublicKey)
		publicKey = pem.EncodeToMemory(&pem.Block{
			Type:  BlockTypeRsaPublicKey,
			Bytes: publicBytes,
		})
	}

	if pkcs == PKCS8 {
		privateBytes, _ := x509.MarshalPKCS8PrivateKey(pri)
		privateKey = pem.EncodeToMemory(&pem.Block{
			Type:  BlockTypePrivateKey,
			Bytes: privateBytes,
		})
		publicBytes, _ := x509.MarshalPKIXPublicKey(&pri.PublicKey)
		publicKey = pem.EncodeToMemory(&pem.Block{
			Type:  BlockTypePublicKey,
			Bytes: publicBytes,
		})
	}

	return
}

// VerifyKeyPair verify key pair matches.
// 验证密钥对是否匹配
func (r rsaKeyPair) VerifyKeyPair(publicKey, privateKey []byte) bool {
	pub, err := r.ExportPublicKey(privateKey)
	if err != nil {
		return false
	}
	if string(publicKey) == string(pub) {
		return true
	}
	return false
}

// IsPublicKey whether is public key.
// 是否是公钥
func (r rsaKeyPair) IsPublicKey(publicKey []byte) bool {
	block, _ := pem.Decode(publicKey)
	return isPublicKey(block)
}

// IsPrivateKey whether is private key.
// 是否是私钥
func (r rsaKeyPair) IsPrivateKey(privateKey []byte) bool {
	block, _ := pem.Decode(privateKey)
	return isPrivateKey(block)
}

func isPublicKey(block *pem.Block) bool {
	if block == nil {
		return false
	}
	if block.Type == BlockTypeRsaPublicKey {
		return true
	}
	if block.Type == BlockTypePublicKey {
		return true
	}
	return false
}

func isPrivateKey(block *pem.Block) bool {
	if block == nil {
		return false
	}
	if block.Type == BlockTypeRsaPrivateKey {
		return true
	}
	if block.Type == BlockTypePrivateKey {
		return true
	}
	return false
}

// FormatPublicKey formats public key, adds header, tail and newline character.
// 格式化公钥，添加头尾和换行符
func (r rsaKeyPair) FormatPublicKey(pkcs keyFormat, publicKey []byte) []byte {
	keyHeader, keyTail := "", ""
	if pkcs == PKCS1 {
		keyHeader = HeaderRsaPublicKey + "\n"
		keyTail = TailRsaPublicKey + "\n"
	}
	if pkcs == PKCS8 {
		keyHeader = HeaderPublicKey + "\n"
		keyTail = TailPublicKey + "\n"
	}
	keyBody := stringSplit(strings.Replace(string(publicKey), "\n", "", -1), 64)
	return []byte(keyHeader + keyBody + keyTail)
}

// FormatPrivateKey formats private key, adds header, tail and newline character
// 格式化私，添加头尾和换行符
func (r rsaKeyPair) FormatPrivateKey(pkcs keyFormat, privateKey []byte) []byte {
	keyHeader, keyTail := "", ""
	if pkcs == PKCS1 {
		keyHeader = HeaderRsaPrivateKey + "\n"
		keyTail = TailRsaPrivateKey + "\n"
	}
	if pkcs == PKCS8 {
		keyHeader = HeaderPrivateKey + "\n"
		keyTail = TailPrivateKey + "\n"
	}
	keyBody := stringSplit(strings.Replace(string(privateKey), "\n", "", -1), 64)
	return []byte(keyHeader + keyBody + keyTail)
}

// ParsePublicKey parses public key.
// 解析公钥
func (r rsaKeyPair) ParsePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, invalidRsaPublicKeyError()
	}
	if block.Type == BlockTypeRsaPublicKey {
		pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, invalidRsaPublicKeyError()
		}
		return pub, nil
	}
	if block.Type == BlockTypePublicKey {
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, invalidRsaPublicKeyError()
		}
		return pub.(*rsa.PublicKey), err
	}
	return nil, invalidRsaPublicKeyError()
}

// ParsePrivateKey parses private key.
// 解析私钥
func (r rsaKeyPair) ParsePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, invalidRSAPrivateKeyError()
	}
	if block.Type == BlockTypeRsaPrivateKey {
		pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, invalidRSAPrivateKeyError()
		}
		return pri, err
	}
	if block.Type == BlockTypePrivateKey {
		pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, invalidRSAPrivateKeyError()
		}
		return pri.(*rsa.PrivateKey), err
	}
	return nil, invalidRSAPrivateKeyError()
}

// ExportPublicKey exports public key from private key.
// 从私钥里导出公钥
func (r rsaKeyPair) ExportPublicKey(privateKey []byte) (publicKey []byte, err error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		err = invalidRSAPrivateKeyError()
		return
	}

	pri, err := r.ParsePrivateKey(privateKey)
	if err != nil {
		return
	}

	blockType, blockBytes := "", []byte("")
	if block.Type == BlockTypeRsaPrivateKey {
		blockType = BlockTypeRsaPublicKey
		blockBytes = x509.MarshalPKCS1PublicKey(&pri.PublicKey)
	}
	if block.Type == BlockTypePrivateKey {
		blockType = BlockTypePublicKey
		blockBytes, err = x509.MarshalPKIXPublicKey(&pri.PublicKey)
	}
	publicKey = pem.EncodeToMemory(&pem.Block{
		Type:  blockType,
		Bytes: blockBytes,
	})
	return
}

// CompressKey compresses key, removes head, tail and newline character.
// 压缩密钥，去掉头尾和换行符
func (r rsaKeyPair) CompressKey(key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	str := strings.Replace(string(key), "\n", "", -1)
	if isPrivateKey(block) {
		str = strings.Replace(str, HeaderRsaPrivateKey, "", -1)
		str = strings.Replace(str, TailRsaPrivateKey, "", -1)
		str = strings.Replace(str, HeaderPrivateKey, "", -1)
		str = strings.Replace(str, TailPrivateKey, "", -1)
	} else if isPublicKey(block) {
		str = strings.Replace(str, HeaderRsaPublicKey, "", -1)
		str = strings.Replace(str, TailRsaPublicKey, "", -1)
		str = strings.Replace(str, HeaderPublicKey, "", -1)
		str = strings.Replace(str, TailPublicKey, "", -1)
	} else {
		return nil, invalidRSAKeyError()
	}
	return []byte(str), nil
}
