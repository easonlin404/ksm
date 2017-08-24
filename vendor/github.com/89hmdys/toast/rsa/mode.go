package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"hash"
	"crypto/sha1"
)

type CipherMode interface {
	Encrypt(plainText []byte, puk *rsa.PublicKey) ([]byte, error)
	Decrypt(cipherText []byte, prk *rsa.PrivateKey) ([]byte, error)
}

type cipherMode int64

type pkcs1v15Cipher cipherMode

func NewPKCS1v15Cipher() CipherMode {
	return new(pkcs1v15Cipher)
}

func (pkcs1v15 *pkcs1v15Cipher) Encrypt(plainText []byte, puk *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, puk, plainText)
}

func (pkcs1v15 *pkcs1v15Cipher) Decrypt(cipherText []byte, prk *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, prk, cipherText)
}

type SignMode interface {
	Sign(src []byte, hash crypto.Hash, prk *rsa.PrivateKey) ([]byte, error)
	Verify(src []byte, sign []byte, hash crypto.Hash, puk *rsa.PublicKey) error
}

type signMode int64

type pkcs1v15Sign signMode

func NewPKCS1v15Sign() SignMode {
	return new(pkcs1v15Sign)
}

func (pkcs1v15 *pkcs1v15Sign) Sign(src []byte, hash crypto.Hash, prk *rsa.PrivateKey) ([]byte, error) {
	if !hash.Available() {
		return nil, errors.New("unsupport hash type")
	}

	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, prk, hash, hashed)
}

func (pkcs1v15 *pkcs1v15Sign) Verify(src []byte, sign []byte, hash crypto.Hash, puk *rsa.PublicKey) error {
	if !hash.Available() {
		return errors.New("unsupport hash type")
	}

	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)

	return rsa.VerifyPKCS1v15(puk, hash, hashed, sign)
}

type oaepCipher struct {
	h hash.Hash
}

func NewOAEPCipher() CipherMode {
	return &oaepCipher{h:sha1.New()}
}

func (oaep *oaepCipher) Encrypt(plainText []byte, puk *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(oaep.h, rand.Reader, puk, plainText, make([]byte, 0))
}

func (oaep *oaepCipher) Decrypt(cipherText []byte, prk *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptOAEP(oaep.h, rand.Reader, prk, cipherText, make([]byte, 0))
}
