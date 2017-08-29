package ksm

import (
	"crypto/rand"
	"encoding/binary"
)

type TLLVBlock struct {
	Tag         uint64
	BlockLength uint32
	ValueLength uint32
	Value       []byte
}

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

func (t *TLLVBlock) Serialize() []byte {
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

	return out
}

type SKR1TLLVBlock struct {
	TLLVBlock
	IV      []byte
	Payload []byte
}

type DecryptedSKR1Payload struct {
	SK             []byte //Session key
	HU             []byte
	R1             []byte
	IntegrityBytes []byte
}

type CkcR1 struct {
	R1 []byte
}

type CkcDataIv struct {
	IV []byte
}

type CkcEncryptedPayload struct {
	Payload []byte
}

type CkcContentKeyDurationBlock struct {
	*TLLVBlock

	LeaseDuration  uint32 // 16-19, The duration of the lease, if any, in seconds.
	RentalDuration uint32 // 20-23, The duration of the rental, if any, in seconds.
	KeyType        uint32 // 24-27,The key type.
	//Reserved       uint32 // Reserved; set to a fixed value of 0x86d34a3a.
	//Padding        []byte // Random values to fill out the TLLV to a multiple of 16 bytes.

}

func newCkcContentKeyDurationBlock(LeaseDuration, RentalDuration uint32) *CkcContentKeyDurationBlock {
	var value []byte

	LeaseDurationOut := make([]byte, 4)
	binary.BigEndian.PutUint32(LeaseDurationOut, LeaseDuration)

	rentalDurationOut := make([]byte, 4)
	binary.BigEndian.PutUint32(rentalDurationOut, RentalDuration)

	keyTypeOut := make([]byte, 4)
	binary.BigEndian.PutUint32(keyTypeOut, Content_Key_valid_for_lease)

	value = append(value, LeaseDurationOut...)
	value = append(value, rentalDurationOut...)
	value = append(value, keyTypeOut...)
	value = append(value, []byte{0x86, 0xd3, 0x4a, 0x3a}...) //Reserved

	tllv := NewTLLVBlock(Tag_Content_Key_Duration, value)

	return &CkcContentKeyDurationBlock{
		TLLVBlock:      tllv,
		LeaseDuration:  LeaseDuration,
		RentalDuration: RentalDuration,
		KeyType:        Content_Key_valid_for_lease,
	}
}

const (
	Tag_SessionKey_R1             = 0x3d1a10b8bffac2ec
	Tag_SessionKey_R1_integrity   = 0xb349d4809e910687
	Tag_AntiReplaySeed            = 0x89c90f12204106b2
	Tag_R2                        = 0x71b5595ac1521133
	Tag_ReturnRequest             = 0x19f9d4e5ab7609cb
	Tag_AssetID                   = 0x1bf7f53f5d5d5a1f
	Tag_TransactionID             = 0x47aa7ad3440577de
	Tag_ProtocolVersionsSupported = 0x67b8fb79ecce1a13
	Tag_ProtocolVersionUsed       = 0x5d81bcbcc7f61703
	Tag_treamingIndicator         = 0xabb0256a31843974
	Tag_MediaPlaybackState        = 0xeb8efdf2b25ab3a0

	Playback_State_ReadyToStart    = 0xf4dee5a2
	Playback_State_PlayingOrPaused = 0xa5d6739e
	Playback_State_Playing         = 0x4f834330
	Playback_State_Halted          = 0x5991bf20
)

const (
	Field_Tag_Length   = 8
	Field_Block_Length = 4
	Field_Value_Length = 4
)

const (
	Tag_Encrypted_CK         = 0x58b38165af0e3d5a
	Tag_R1                   = 0xea74c4645d5efee9
	Tag_Content_Key_Duration = 0x47acf6a418cd091a
	Tag_HDCP_Enforcement     = 0x2e52f1530d8ddb4a

	Content_Key_valid_for_lease  = 0x1a4bde7e
	Content_key_valid_for_rental = 0x3dfe45a0
	Content_key_valid_for_both   = 0x27b59bde
)
