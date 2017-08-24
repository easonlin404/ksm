package cipher

import . "crypto/cipher"


/*
介绍:CipherMode为不同的工作模式提供了统一的接口来设置填充方式，创建Cipher。
    对于流模式，SetPadding是个空方法，不起任何作用。

作者:Alex
版本:release-1.1
*/
type CipherMode  interface {
	/*
	设置填充方式,仅块模式需要填充。

	作者:Alex
	版本:release-1.1
	*/
	SetPadding(padding Padding) CipherMode

	/*
	创建Cipher

	作者:Alex
	版本:release-1.1
	*/
	Cipher(block Block, iv []byte) Cipher
}

type cipherMode struct {
	padding Padding
}

func (c *cipherMode) SetPadding(padding Padding) CipherMode {
	return c
}

func (c *cipherMode) Cipher(block Block, iv []byte) Cipher {
	return nil
}

type ecbCipherModel cipherMode

func NewECBMode() CipherMode {
	return &ecbCipherModel{padding:NewPKCS57Padding() }
}

func (ecb *ecbCipherModel) SetPadding(padding Padding) CipherMode {
	ecb.padding = padding
	return ecb
}

func (ecb *ecbCipherModel) Cipher(block Block, iv []byte) Cipher {
	encrypter := NewECBEncrypter(block)
	decrypter := NewECBDecrypter(block)
	return NewBlockCipher(ecb.padding, encrypter, decrypter)
}

type cbcCipherModel cipherMode

func NewCBCMode() CipherMode {
	return &cbcCipherModel{padding:NewPKCS57Padding()}
}

func (cbc *cbcCipherModel) SetPadding(padding Padding) CipherMode {
	cbc.padding = padding
	return cbc
}

func (cbc *cbcCipherModel) Cipher(block Block, iv []byte) Cipher {
	encrypter := NewCBCEncrypter(block, iv)
	decrypter := NewCBCDecrypter(block, iv)
	return NewBlockCipher(cbc.padding, encrypter, decrypter)
}

type cfbCipherModel cipherMode

func NewCFBMode() CipherMode {
	return &ofbCipherModel{}
}

func (cfb *cfbCipherModel) Cipher(block Block, iv []byte) Cipher {
	encrypter := NewCFBEncrypter(block, iv)
	decrypter := NewCFBDecrypter(block, iv)
	return NewStreamCipher(encrypter, decrypter)
}

type ofbCipherModel struct {
	cipherMode
}

func NewOFBMode() CipherMode {
	return &ofbCipherModel{}
}

func (ofb *ofbCipherModel) Cipher(block Block, iv []byte) Cipher {
	encrypter := NewOFB(block, iv)
	decrypter := NewOFB(block, iv)
	return NewStreamCipher(encrypter, decrypter)
}

type ctrCipherModel struct {
	cipherMode
}

func NewCTRMode() CipherMode {
	return &ctrCipherModel{}
}

func (ctr *ctrCipherModel) Cipher(block Block, iv []byte) Cipher {
	encrypter := NewCTR(block, iv)
	decrypter := NewCTR(block, iv)
	return NewStreamCipher(encrypter, decrypter)
}