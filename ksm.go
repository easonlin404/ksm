package ksm

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	_ "fmt"
	"github.com/easonlin404/ksm/aes"
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

func ParseSPCV1(playback []byte, pem []byte) (*SPCContainer, error) {
	spcContainer := parseSPCContainer(playback)

	spck, err := decryptSPCK(pem, spcContainer.EncryptedAesKey)
	if err != nil {
		return nil, err
	}

	spcpayload, err := decryptSPCpayload(spcContainer, spck)
	if err != nil {
		return nil, err
	}
	fmt.Println(spcpayload)

	//TODO:
	printDebugSPC(spcContainer)

	return spcContainer, nil
}

func parseSPCContainer(playback []byte) *SPCContainer {
	spcContainer := &SPCContainer{}
	spcContainer.Version = binary.BigEndian.Uint32(playback[0:4])
	spcContainer.AesKeyIV = playback[8:24]
	spcContainer.EncryptedAesKey = playback[24:152]
	spcContainer.CertificateHash = playback[152:172]
	spcContainer.SPCPlayloadLength = binary.BigEndian.Uint32(playback[172:176])
	spcContainer.SPCPlayload = playback[176 : 176+spcContainer.SPCPlayloadLength]

	return spcContainer
}

func parseTLLVs(spcpayload []byte)map[string]TLLVBlock {
	var m map[string]TLLVBlock

	m = make(map[string]TLLVBlock)

	return m
}

func printDebugSPC(spcContainer *SPCContainer) {
	fmt.Println("========================= Begin SPC Data ===============================")
	fmt.Println("SPC container size -")
	fmt.Println(spcContainer.SPCPlayloadLength)

	fmt.Println("SPC Encryption Key -")
	fmt.Println(hex.EncodeToString(spcContainer.EncryptedAesKey))
	fmt.Println("SPC Encryption IV -")
	fmt.Println(hex.EncodeToString(spcContainer.AesKeyIV))
	fmt.Println("================ SPC TLLV List ================")
	//TODO:
	fmt.Println("[SK ... R1] Integrity Tag --")
	fmt.Println("=========================== End SPC Data =================================")

}

// SPCK = RSA_OAEP d([SPCK])Prv where
// [SPCK] represents the value of SPC message bytes 24-151. Prv represents the server's private key.
func decryptSPCK(pkPem, enSpck []byte) ([]byte, error) {
	if len(enSpck) != 128 {
		return nil, errors.New("Wrong [SPCK] length, must be 128")
	}
	return rsa.OAEPPDecrypt(pkPem, enSpck)
}

// SPC payload = AES_CBCIV d([SPC data])SPCK where
// [SPC data] represents the remaining SPC message bytes beginning at byte 176 (175 + the value of
// SPC message bytes 172-175).
// IV represents the value of SPC message bytes 8-23.
func decryptSPCpayload(spcContainer *SPCContainer, spck []byte) ([]byte, error) {
	spcPayload, err := aes.Decrypt(spck, spcContainer.AesKeyIV, spcContainer.SPCPlayload)
	return spcPayload, err
}
