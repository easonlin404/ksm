package aes

import (
	"github.com/89hmdys/toast/crypto"
	"github.com/89hmdys/toast/cipher"
)

func EncryptWithECB(key []byte, plainText []byte) ([]byte, error) {
	mode := cipher.NewECBMode()
	cipher, err := crypto.NewAESWith(key, mode)
	if err != nil {
		return nil,err
	}

	ciphertext := cipher.Encrypt(plainText)
	return ciphertext,nil
}