package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// Encrypt encrypts the given message with RSA-OAEP.
// Need a DER encoded public key, These values are
// typically found in PEM blocks with "BEGIN PUBLIC KEY".
func Encrypt(publicKey, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)

	if block == nil {
		return nil, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	//fmt.Println("Modulus : ", pub.N.String())
	//fmt.Println(">>> ", pub.N)
	//fmt.Printf("Modulus(Hex) : %X\n", pub.N)
	//fmt.Println("Public Exponent : ", pub.E)
	//return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, origData, nil)
}

func EncryptByCert(pemCertificate, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(pemCertificate)
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	pub := cert.PublicKey.(*rsa.PublicKey)

	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, origData, nil)
}

func Decrypt(privateKey, ciphertext []byte) ([]byte, error) {
	fmt.Println(len(ciphertext))
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) // ASN.1 PKCS#1 DER encoded form.
	if err != nil {
		return nil, err
	}
	//return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, nil)
}
