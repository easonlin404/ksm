package rsa

import "fmt"

type Padding interface {
	Padding(src []byte) [][]byte
}

type padding struct {
	lessThanModulus int
}

func (padding *padding) Padding(src []byte) [][]byte {
	return grouping(src, padding.lessThanModulus)
}

func NewPKCS1Padding(modulus int) Padding {
	paddingLen := 11
	return newPadding(modulus - paddingLen)
}

func NewOAEPPadding(modulus int) Padding {
	paddingLen := 41
	return newPadding(modulus - paddingLen)
}

func NewNoPadding(modulus int) Padding {
	return newPadding(modulus)
}

func newPadding(lessThanModulus int) Padding {
	return &padding{lessThanModulus: lessThanModulus}
}

/*数据太长的时候，要按照秘钥的长度对数据进行分组*/
func grouping(src []byte, size int) [][]byte {

	var groups [][]byte

	fmt.Println(size)

	srcSize := len(src)
	if srcSize <= size {
		groups = append(groups, src)
	} else {
		for len(src) != 0 {
			if len(src) <= size {
				groups = append(groups, src)
				break
			} else {
				v := src[:size]
				groups = append(groups, v)
				src = src[size:]
			}
		}
	}

	return groups
}
