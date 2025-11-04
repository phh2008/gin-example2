package openssl

const (
	BlockTypeRsaPrivateKey = "RSA PRIVATE KEY"
	BlockTypePrivateKey    = "PRIVATE KEY"

	BlockTypeRsaPublicKey = "RSA PUBLIC KEY"
	BlockTypePublicKey    = "PUBLIC KEY"
)

const (
	HeaderRsaPrivateKey = "-----BEGIN RSA PRIVATE KEY-----"
	TailRsaPrivateKey   = "-----END RSA PRIVATE KEY-----"
	HeaderPrivateKey    = "-----BEGIN PRIVATE KEY-----"
	TailPrivateKey      = "-----END PRIVATE KEY-----"

	HeaderRsaPublicKey = "-----BEGIN RSA PUBLIC KEY-----"
	TailRsaPublicKey   = "-----END RSA PUBLIC KEY-----"
	HeaderPublicKey    = "-----BEGIN PUBLIC KEY-----"
	TailPublicKey      = "-----END PUBLIC KEY-----"
)
