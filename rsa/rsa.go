package rsa

import (
	"encoding/pem"
	"crypto/rand"
	"crypto/x509"
	"crypto/rsa"
	"errors"
)


func Encrypt(publicKey,origData []byte) ([]byte, error) {
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
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

//func EncryptByCert(cert,origData []byte) ([]byte, error) {
//	block, _ := pem.Decode(cert)
//	if block == nil {
//		panic("failed to parse certificate PEM")
//	}
//	cert, err := x509.ParseCertificate(block.Bytes)
//	if err != nil {
//		panic("failed to parse certificate: " + err.Error())
//	}
//	rsa.en
//	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
//}


func Decrypt(privateKey,ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	//label := []byte("orders")
	//return rsa.DecryptOAEP(sha256.New(),rand.Reader,priv,ciphertext,label)
}

