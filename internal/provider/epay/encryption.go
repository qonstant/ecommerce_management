package epay

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

const PublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAqoATnGMtByQojuoYFKEx
vEqszShV2vk6chCJFx0/vmSHBqcCTazhJqBmYU9gyM/TjVWLsjFAbd4nvCxIGpqF
g3J7UJccfODKibyfSUSqGsAJE1MJh3EaJivkd85/FkZkv3zBeT/193NmakNs0+T+
PUMmSdAPSnfUWi2KSIp48mSA38CbMvOwndkKNEeqCoIQn/fApfZ8MWIEFVd3gpfs
Ve0zYhSjvTOHPD0/7TOdcQyArxLZY0yS7m32rUOibuO7EhGNQL/bC73ZbuS5nXhr
a03nNIW3FfSJUTJBjVWDZRoNk9gm4pOimAeb0IiqnmlTPOvkqHYsOEjQ8KJAFlGO
1igelk1+dA5ZiY6r0YExc1KnW7UsnGk6nr7cgOR2po/sa4kctiKLqlGA35ILmUBQ
Yb6iReCQkggXMOvmP6p+4wEt1B7V8UJxzFZcQZ5QSRIk3o3pVrfY0gksidl0Xt5m
ft+E6a77ZQKG4TOQS9Ly1mIJ2qqaWqCWglVMWFiFCx9dXTN0RMli1T0rs1gA2jsP
z2/HiyY8EUp6t4Ufc8VbJYG9vt24UTwYgu+qDEBjggm5YKVCxjCvhJWwh9LaL9Uu
K46Apgr5wgEyMIJZRO7RxkjKkJI29FAP3wEs9y+/3qsjH3chFzdX0/+6lA+9lePK
PX0Z5SPexWRiQp9bND4iZRcCAwEAAQ==
-----END PUBLIC KEY-----`

// Encrypts data using RSA public key
func EncryptWithPublicKey(data []byte, publicKeyPEM string) (string, error) {
	// Decode the PEM public key
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing public key: %v", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error asserting public key type")
	}

	// Encrypt the data
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, data)
	if err != nil {
		return "", fmt.Errorf("error encrypting data: %v", err)
	}

	// Encode the encrypted data to base64
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}
