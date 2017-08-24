package cipher

import "bytes"

/*
介绍:Padding接口为各种填充方式提供了统一的接口来填充/还原数据。
作者:Alex
版本:release-1.1
*/
type Padding interface {
	/*
	介绍:根据块大小填充待加密数据
        作者:Alex
        版本:release-1.1
	*/
	Padding(src []byte, blockSize int) []byte
	/*
	介绍:从解密后的数据中取出填充的数据，还原原文
        作者:Alex
        版本:release-1.1
	*/
	UnPadding(src []byte) []byte
}

type padding struct {}

type pkcs57Padding  padding

/*
介绍:创建PKCS5/7填充模式
作者:Alex
版本:release-1.1
*/
func NewPKCS57Padding() Padding {
	return &pkcs57Padding{}
}

func (p *pkcs57Padding) Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func (p *pkcs57Padding) UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length - 1])
	return src[:(length - unpadding)]
}
