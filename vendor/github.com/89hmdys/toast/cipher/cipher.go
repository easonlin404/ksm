package cipher

import (
	. "crypto/cipher"
)


/*
介绍:Cipher提供了统一的接口对数据进行加密/解密操作.

作者:Alex
版本:release-1.1
*/
type Cipher interface {
	/*
	介绍:加密数据
	作者:Alex
        版本:release1.1
	*/
	Encrypt(src []byte) []byte
	/*
	介绍:解密数据
	作者:Alex
        版本:release1.1
	*/
	Decrypt(src []byte) []byte
}

/*
介绍:新建块加密
作者:Alex
版本:release1.1
*/
func NewBlockCipher(padding Padding, encrypt, decrypt BlockMode) Cipher {
	return &blockCipher{
		encrypt:   encrypt,
		decrypt:   decrypt,
		padding:   padding}
}

type blockCipher struct {
	padding Padding
	encrypt BlockMode
	decrypt BlockMode
}

func (blockCipher *blockCipher) Encrypt(plaintext []byte) []byte {
	//TODO: modify By Eason ECB nopadding
	//plaintext = blockCipher.padding.Padding(plaintext, blockCipher.encrypt.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	blockCipher.encrypt.CryptBlocks(ciphertext, plaintext)
	return ciphertext
}

func (blockCipher *blockCipher) Decrypt(ciphertext []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	blockCipher.decrypt.CryptBlocks(plaintext, ciphertext)
	plaintext = blockCipher.padding.UnPadding(plaintext)
	return plaintext
}

/*
介绍:新建流加密
作者:Alex
版本:release1.1
*/
func NewStreamCipher(encrypt Stream, decrypt Stream) Cipher {
	return &streamCipher{
		encrypt: encrypt,
		decrypt: decrypt}
}

type streamCipher struct {
	encrypt Stream
	decrypt Stream
}

func (streamCipher *streamCipher) Encrypt(plaintext []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	streamCipher.encrypt.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}
func (streamCipher *streamCipher) Decrypt(ciphertext []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	streamCipher.decrypt.XORKeyStream(plaintext, ciphertext)
	return plaintext
}
