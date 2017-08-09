package ksm

import (
	"encoding/binary"
	"github.com/easonlin404/ksm/rsa"
)

type SPCContainer struct {
	Version           uint32
	Reserved          []byte
	AesKeyIV          []byte //16
	EncryptedAesKey   []byte //128
	CertificateHash   []byte //20
	SPCPlayload       []byte //TODO: struct
	SPCPlayloadLength uint32
}



// This function will compute the content key context returned to client by the SKDServer library.
//       incoming server playback context (SPC message)
func GenCKC(playback []byte) error {

	return nil
}

func ParseSPCContainer(playback []byte,reader rsa.Reader) (*SPCContainer, error) {
	spcContainer := &SPCContainer{}

	spcContainer.Version = binary.BigEndian.Uint32(playback[0:4])
	spcContainer.AesKeyIV = playback[8:24]

	//TODO: get RSA private key
	pem,err:=reader.ReadPem()
	if err != nil {
		return spcContainer, err
	}

	encryptedAesKey, err := rsa.Decrypt(pem, playback[24:152])
	if err != nil {
		return spcContainer, err
	}
	spcContainer.EncryptedAesKey = encryptedAesKey

	return spcContainer, nil
}
