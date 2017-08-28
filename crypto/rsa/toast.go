package rsa

import (
	"github.com/89hmdys/toast/crypto"
	"github.com/89hmdys/toast/rsa"
)

//RSA/ECB/OAEPPadding
func OAEPPDecrypt(pub, pri string, ciphertext []byte) ([]byte, error) {
	key, err := rsa.LoadKeyFromPEMByte(
		[]byte(pub),
		[]byte(pri),
		rsa.ParsePKCS1KeyByCert)

	if err != nil {
		panic(err)
	}

	padding := rsa.NewOAEPPadding(key.Modulus())
	cipherMode := rsa.NewOAEPCipher()
	signMode := rsa.NewPKCS1v15Sign()

	cipher, err := crypto.NewRSAWith(key, padding, cipherMode, signMode)
	if err != nil {
		panic(err)
	}

	deT, err := cipher.Decrypt(ciphertext)
	if err != nil {
		panic(err)
	}
	return deT, err
}
