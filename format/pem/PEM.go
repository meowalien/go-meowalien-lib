package pem

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func PEMBytesToPublicKey(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pp := pub.(type) {
	case *rsa.PublicKey:
		return pp, nil
	default:
		return nil, errors.New("Key type is not RSA")
	}
}

func PublicKeyToPEMBytes(key *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	return pubBytes, nil
}
