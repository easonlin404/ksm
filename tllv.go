package ksm

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
)

// TLLVBlock represents a TLLV block structure.
type TLLVBlock struct {
	Tag         uint64
	BlockLength uint32
	ValueLength uint32 //The number of bytes in the value field. This number may be any amount, including 0x0000
	Value       []byte
}

// NewTLLVBlock creates a new TLLVBlock object using the specified tag and value.
func NewTLLVBlock(tag uint64, value []byte) *TLLVBlock {
	valueLen := uint32(len(value))
	paddingSize := 32 - valueLen%16 // Extend to nearest 16 bytes + extra 16 bytes
	blockLen := valueLen + paddingSize

	return &TLLVBlock{
		Tag:         tag,
		BlockLength: blockLen,
		ValueLength: valueLen,
		Value:       value,
	}
}

// Serialize returns serialize byte array.
func (t *TLLVBlock) Serialize() ([]byte, error) {
	if err := t.check(); err != nil {
		return nil, err
	}

	var out []byte

	tagOut := make([]byte, 8)
	blockLenOut := make([]byte, 4)
	valueLenOut := make([]byte, 4)

	valueLen := uint32(len(t.Value))
	paddingLen := 32 - valueLen%16 // Extend to nearest 16 bytes + extra 16 bytes
	blockLen := valueLen + paddingLen

	paddingOut := make([]byte, paddingLen)
	rand.Read(paddingOut)

	binary.BigEndian.PutUint64(tagOut, t.Tag)
	binary.BigEndian.PutUint32(blockLenOut, blockLen)
	binary.BigEndian.PutUint32(valueLenOut, valueLen)

	out = append(out, tagOut...)
	out = append(out, blockLenOut...)
	out = append(out, valueLenOut...)
	out = append(out, t.Value...)
	out = append(out, paddingOut...)

	return out, nil
}

func (t *TLLVBlock) check() error {
	if t.Tag == 0 {
		return errors.New("tag not found")
	}
	if len(t.Value) == 0 {
		fmt.Printf("tag: %x :value not found\n", t.Tag)
		fmt.Printf("tag.ValueLength: %x \n", t.ValueLength)

		//return fmt.Errorf("tag: %x :value not found", t.Tag)
	}
	return nil
}

// SKR1TLLVBlock represents a SKR1 TLLV block structure.
type SKR1TLLVBlock struct {
	TLLVBlock
	IV      []byte
	Payload []byte
}

// DecryptedSKR1Payload represents a decrypted SKR1 payload structure.
type DecryptedSKR1Payload struct {
	SK             []byte //Session key
	HU             []byte
	R1             []byte
	IntegrityBytes []byte
}

// CkcR1 represents a ckcR1 structure.
type CkcR1 struct {
	R1 []byte
}

// CkcDataIv represents a ckc data iv structure.
type CkcDataIv struct {
	IV []byte
}

//CkcEncryptedPayload represents a ckc encrypted payload structure.
type CkcEncryptedPayload struct {
	Payload []byte
}

// CkcContentKeyDurationBlock represents a ckc content key duration block structure.
type CkcContentKeyDurationBlock struct {
	*TLLVBlock

	LeaseDuration  uint32 // 16-19, The duration of the lease, if any, in seconds.
	RentalDuration uint32 // 20-23, The duration of the rental, if any, in seconds.
	KeyType        uint32 // 24-27,The key type.
	//Reserved       uint32 // Reserved; set to a fixed value of 0x86d34a3a.
	//Padding        []byte // Random values to fill out the TLLV to a multiple of 16 bytes.

}

// NewCkcContentKeyDurationBlock creates a new a ckc content key duration block object using the specified lease duration and rental duration.
func NewCkcContentKeyDurationBlock(LeaseDuration, RentalDuration uint32) *CkcContentKeyDurationBlock {
	var value []byte

	LeaseDurationOut := make([]byte, 4)
	binary.BigEndian.PutUint32(LeaseDurationOut, LeaseDuration)

	rentalDurationOut := make([]byte, 4)
	binary.BigEndian.PutUint32(rentalDurationOut, RentalDuration)

	keyTypeOut := make([]byte, 4)

	//binary.BigEndian.PutUint32(keyTypeOut, contentKeyValidForLease)
	binary.BigEndian.PutUint32(keyTypeOut, contentKeyPersisted)

	value = append(value, LeaseDurationOut...)
	value = append(value, rentalDurationOut...)
	value = append(value, keyTypeOut...)
	value = append(value, []byte{0x86, 0xd3, 0x4a, 0x3a}...) //Reserved; set to a fixed value of 0x86d34a3a.

	tllv := NewTLLVBlock(tagContentKeyDuration, value)

	return &CkcContentKeyDurationBlock{
		TLLVBlock:      tllv,
		LeaseDuration:  LeaseDuration,
		RentalDuration: RentalDuration,
		KeyType:        contentKeyValidForLease, //TODO: KeyType does not use
	}
}

const (
	tagSessionKeyR1              = 0x3d1a10b8bffac2ec
	tagSessionKeyR1Integrity     = 0xb349d4809e910687
	tagAntiReplaySeed            = 0x89c90f12204106b2
	tagR2                        = 0x71b5595ac1521133
	tagReturnRequest             = 0x19f9d4e5ab7609cb
	tagAssetID                   = 0x1bf7f53f5d5d5a1f
	tagTransactionID             = 0x47aa7ad3440577de
	tagProtocolVersionsSupported = 0x67b8fb79ecce1a13
	tagProtocolVersionUsed       = 0x5d81bcbcc7f61703
	tagTreamingIndicator         = 0xabb0256a31843974
	tagMediaPlaybackState        = 0xeb8efdf2b25ab3a0

	playbackStateReadyToStart    = 0xf4dee5a2
	playbackStatePlayingOrPaused = 0xa5d6739e
	playbackStatePlaying         = 0x4f834330
	playbackStateHalted          = 0x5991bf20
)

const (
	fieldTagLength   = 8
	fieldBlockLength = 4
	fieldValueLength = 4
)

const (
	tagEncryptedCk        = 0x58b38165af0e3d5a
	tagR1                 = 0xea74c4645d5efee9
	tagContentKeyDuration = 0x47acf6a418cd091a
	tagHdcpEnforcement    = 0x2e52f1530d8ddb4a

	contentKeyValidForLease  = 0x1a4bde7e //Content key valid for lease only
	contentKeyValidForRental = 0x3dfe45a0 //Content key valid for rental only
	contentKeyValidForBoth   = 0x27b59bde //Content key valid for both lease and rental
)

const (
	//Offline
	contentKeyPersisted            = 0x3df2d9fb //Content key can be persisted with unlimited validity duration
	contentKeyPersistedWithlimited = 0x18f06048 //Content key can be persisted, and it’s validity duration is limited to the “Rental Duration” value
)
