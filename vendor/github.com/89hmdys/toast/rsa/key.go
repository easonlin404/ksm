package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"strings"
)

type Key interface {
	PublicKey() *rsa.PublicKey
	PrivateKey() *rsa.PrivateKey
	Modulus() int
}

func ParsePKCS8Key(publicKey, privateKey []byte) (Key, error) {
	puk, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	prk, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return &key{publicKey: puk.(*rsa.PublicKey), privateKey: prk.(*rsa.PrivateKey)}, nil
}

func ParsePKCS1Key(publicKey, privateKey []byte) (Key, error) {
	puk, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	prk, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return &key{publicKey: puk.(*rsa.PublicKey), privateKey: prk}, nil
}

func ParsePKCS1KeyByCert(publicKey, privateKey []byte) (Key, error) {
	puk, err := x509.ParseCertificate(publicKey)
	if err != nil {
		return nil, err
	}
	prk, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return &key{publicKey: puk.PublicKey.(*rsa.PublicKey), privateKey: prk}, nil
}

func LoadKeyFromPEMFile(publicKeyFilePath, privateKeyFilePath string, ParseKey func([]byte, []byte) (Key, error)) (Key, error) {

	//TODO 断言如果入参为"" ，则直接报错

	publicKeyFilePath = strings.TrimSpace(publicKeyFilePath)

	pukBytes, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, err
	}

	puk, _ := pem.Decode(pukBytes)
	if puk == nil {
		return nil, errors.New("publicKey is not pem formate")
	}

	privateKeyFilePath = strings.TrimSpace(privateKeyFilePath)

	prkBytes, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}

	prk, _ := pem.Decode(prkBytes)
	if prk == nil {
		return nil, errors.New("privateKey is not pem formate")
	}

	return ParseKey(puk.Bytes, prk.Bytes)
}


//Eason Lin
func LoadKeyFromPEMByte(pukBytes, prkBytes []byte, ParseKey func([]byte, []byte) (Key, error)) (Key, error) {
	puk, _ := pem.Decode(pukBytes)
	if puk == nil {
		return nil, errors.New("publicKey is not pem formate")
	}

	prk, _ := pem.Decode(prkBytes)
	if prk == nil {
		return nil, errors.New("privateKey is not pem formate")
	}

	return ParseKey(puk.Bytes, prk.Bytes)
}


func LoadKeyFromDerFile(publicKeyFilePath, privateKeyFilePath string, ParseKey func([]byte, []byte) (Key, error)) (Key, error) {

	publicKeyFilePath = strings.TrimSpace(publicKeyFilePath)

	pukBytes, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, err
	}

	privateKeyFilePath = strings.TrimSpace(privateKeyFilePath)

	prkBytes, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}

	return ParseKey(pukBytes, prkBytes)
}



type key struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (key *key) Modulus() int {
	return len(key.publicKey.N.Bytes())
}

func (key *key) PublicKey() *rsa.PublicKey {
	return key.publicKey
}

func (key *key) PrivateKey() *rsa.PrivateKey {
	return key.privateKey
}
