package ksm

import (
	"crypto/rand"
	mathRand "math/rand"
)

// ContentKey is a interface that fetch asset content key and duration.
type ContentKey interface {
	FetchContentKey(assetId []byte) ([]byte, []byte, error)
	FetchContentKeyDuration(assetId []byte) (*CkcContentKeyDurationBlock, error)
}

var (
	_ ContentKey = RandomContentKey{}
)

// CKCPayload is a object that implements ContentKey interface.
type RandomContentKey struct {
}

func (RandomContentKey) FetchContentKey(assetId []byte) ([]byte, []byte, error) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	rand.Read(key)
	rand.Read(iv)
	return key, iv, nil
}

// FetchContentKeyDuration returns CkcContentKeyDurationBlock for the given assetId.
func (RandomContentKey) FetchContentKeyDuration(assetId []byte) (*CkcContentKeyDurationBlock, error) {

	LeaseDuration := mathRand.Uint32()  // The duration of the lease, if any, in seconds.
	RentalDuration := mathRand.Uint32() // The duration of the rental, if any, in seconds.

	return NewCkcContentKeyDurationBlock(LeaseDuration, RentalDuration), nil
}

// CKCPayload is a object that store ckc payload.
type CKCPayload struct {
	SK             []byte //Session key
	HU             []byte
	R1             []byte
	IntegrityBytes []byte
}
