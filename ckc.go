package ksm

import "crypto/rand"

type ContentKey interface {
	FetchContentKeyAndIV(assetId []byte) ([]byte,[]byte, error)
}

var (
	_ ContentKey = RandomContentKey{}
)

type RandomContentKey struct {
}

func (RandomContentKey) FetchContentKeyAndIV(assetId []byte) ([]byte,[]byte, error) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	rand.Read(key)
	rand.Read(iv)

	return key,iv,nil
}

type CKCPayload struct {
	SK             []byte //Session key
	HU             []byte
	R1             []byte
	IntegrityBytes []byte
}
