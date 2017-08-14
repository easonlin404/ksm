package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func Encrypt(key []byte, iv []byte, plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	//content := PKCS5Padding(plainText, block.BlockSize())
	//crypted := make([]byte, len(content))
	cbc.CryptBlocks(plainText, plainText)
	return plainText, nil
}

func Decrypt(key []byte, iv []byte, crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(crypted, crypted)
	//crypted = PKCS5UnPadding(crypted)
	return crypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
