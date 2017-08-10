package ksm

import (
	"encoding/binary"
	"errors"
	_ "fmt"
	"github.com/easonlin404/ksm/rsa"
	"fmt"
	"encoding/hex"
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

func ParseSPCContainer(playback []byte, reader rsa.Reader) (*SPCContainer, error) {
	fmt.Println("========================= Begin SPC Data ===============================")

	spcContainer := &SPCContainer{}

	spcContainer.Version = binary.BigEndian.Uint32(playback[0:4])
	spcContainer.AesKeyIV = playback[8:24]
	spck := playback[24:152]
	spcContainer.EncryptedAesKey = spck


	//encryptedAesKey, err := rsa.Decrypt(pem, playback[24:152])

	//fmt.Println("[SPCK]:" + hex.EncodeToString(spck))
	//fmt.Println(len(hex.EncodeToString(spck)))
	//pem, err := reader.ReadPem()
	//if err != nil {
	//	return spcContainer, err
	//}
	//encryptedAesKey, err := decryptSPCK(pem, spck) //TODO: not this





	fmt.Println("SPC Encryption Key -")
	fmt.Println(hex.EncodeToString(spck))
	fmt.Println("SPC Encryption IV -")
	fmt.Println(hex.EncodeToString(spcContainer.AesKeyIV))
	fmt.Println("================ SPC TLLV List ================")
	//TODO:
	return spcContainer, nil
}

func decryptSPCK(pkPem, enSpck []byte) ([]byte, error) {
	if len(enSpck) != 128 {
		return nil, errors.New("Wrong [SPCK] length, must be 128")
	}
	return rsa.OAEPPDecrypt(pkPem, enSpck)
}
