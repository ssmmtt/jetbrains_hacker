package algo

import (
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

// The following code is copied from the standard library
var hashPrefixes = map[string][]byte{
	"SHA-256": {0x30, 0x31, 0x30, 0x0d, 0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01, 0x05, 0x00, 0x04, 0x20},
}

var ErrMessageTooLong = errors.New("crypto/rsa: message too long for RSA key size")

func pkcs1v15ConstructEM(pub *rsa.PublicKey, hash string, hashed []byte) ([]byte, error) {
	// Special case: "" is used to indicate that the data is signed directly.
	var prefix []byte
	if hash != "" {
		var ok bool
		prefix, ok = hashPrefixes[hash]
		if !ok {
			return nil, errors.New("crypto/rsa: unsupported hash function")
		}
	}

	// EM = 0x00 || 0x01 || PS || 0x00 || T
	k := pub.Size()
	if k < len(prefix)+len(hashed)+2+8+1 {
		return nil, ErrMessageTooLong
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-len(prefix)-len(hashed)-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-len(prefix)-len(hashed):], prefix)
	copy(em[k-len(hashed):], hashed)
	return em, nil
}

func GetEM(pub *rsa.PublicKey, tbsCertificate []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(tbsCertificate)
	signed := h.Sum(nil)
	return pkcs1v15ConstructEM(pub, "SHA-256", signed)
}
