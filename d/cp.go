package d

import (
	"crypto/sha1"
	"errors"
	"github.com/easonlin404/ksm/aes"
)

const PRIME = uint32(813416437)
const NB_RD = 16

type CP_D_Function struct {
}

func (d CP_D_Function) Compute(R2 []byte, R2_sz uint32, ask []byte) ([]byte, error) {
	hh := d.ComputeHashValue(R2, R2_sz)
	DASk, err := aes.EncryptWithECB(ask, hh)

	if err != nil {
		return nil, err
	}
	/*
		* @param[input]      data to encrypt or decrypt
		* @param[output]     encrypted or decrypted data
		* @param[inputSize]  size of data to encrypt
		* @param[key]        key used to do the encryption or decryption
		* @param[iv]         iv used to the encryption or decryption
		* @param[opType]     kSKDServerAESEncrypt, or kSKDServerAESDecrypt
		* @param[opMode]     kSKDServerAES_CBC, or kSKDServerAES_ECB
		status = SKDServerAESEncryptDecrypt(tmp, DASk, 16, (UInt8*)ASk, NULL, kSKDServerAESEncrypt, kSKDServerAES_ECB);
	*/
	//
	if len(DASk) != 16 {
		errors.New("DASk key length doesn't equal 16")
	}

	return DASk, nil

}

func (d CP_D_Function) ComputeHashValue(R2 []byte, R2_sz uint32) []byte {
	//var R2_sz := 21

	//var tmp [20]uint8
	var pad []byte
	pad = make([]byte, 64, 64)
	var MBlock [14]uint32
	var r uint32
	var P uint32

	var i uint32

	P = PRIME

	if len(R2) == 0 {
		errors.New("R2 block doesn't exist.")
	}

	/* Padding until a multiple of 56B */
	for i = 0; i < R2_sz; i++ {
		pad[i] = R2[i]
	}

	pad[R2_sz] = 0x80
	for i = R2_sz + 1; i < 56; i++ {
		pad[i] = 0
	}

	/* Create 14 32b values */
	for i = 0; i < 14; i++ {
		MBlock[i] = uint32((pad[4*i] << 24) ^ (pad[4*i+1] << 16) ^ (pad[4*i+2] << 8) ^ (pad[4*i+3]))
	}

	/* Reunify into 2 32 bits values */
	for i = 1; i < 7; i++ {
		MBlock[0] += MBlock[i]
	}

	MBlock[1] = 0
	for i = 0; i < 7; i++ {
		MBlock[1] += MBlock[i+7]
	}

	// Apply the function
	// This block is the C_r function specified in the D Function specification document
	for i = 0; i < 2; i++ {
		for r = 0; r < NB_RD; r++ {
			if MBlock[i]&1 == 1 {
				MBlock[i] >>= 1
			} else {
				MBlock[i] = (3*MBlock[i] + 1) % P
			}
		}
	}

	/* append to M */
	for i = 0; i < 4; i++ {
		pad[56+i] = uint8(MBlock[0] >> (8 * i))
		pad[60+i] = uint8(MBlock[1] >> (8 * i))
	}
	//pad have to 64 size

	h := sha1.New()
	h.Write(pad)

	hh := h.Sum(nil)

	//if len(hh)!=20{
	//	panic("hash value must length 20  expected.")//or 20?
	//}

	return hh[0:16]

}
