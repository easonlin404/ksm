package ksm

import "encoding/hex"

type SPCContainer struct {
	Version           uint32
	AesKeyIV          []byte //16
	EncryptedAesKey   []byte //128
	CertificateHash   []byte //20
	SPCPlayload       []byte //TODO: struct
	SPCPlayloadLength uint32
}

// This function will compute the content key context(hex) returned to client by the SKDServer library.
//       incoming server playback context (SPC message)
func GenCKC(playbackHex string) error {
	b, err := hex.DecodeString(playbackHex)
	if err != nil {
		return err
	}
	println(len(b))

	return nil
}
