package ksm

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"

	"github.com/easonlin404/ksm/aes"
	"github.com/easonlin404/ksm/d"
	"github.com/easonlin404/ksm/rsa"
)

type SPCContainer struct {
	Version           uint32
	Reserved          []byte
	AesKeyIV          []byte
	EncryptedAesKey   []byte
	CertificateHash   []byte
	SPCPlayload       []byte
	SPCPlayloadLength uint32

	TTLVS map[uint64]TLLVBlock
}

// This function will compute the content key context returned to client by the SKDServer library.
//       incoming server playback context (SPC message)
func GenCKC(playback []byte) error {
	pem := []byte{} //TODO: server pk
	spcv1, err := ParseSPCV1(playback, pem)
	if err != nil {
		return err
	}

	ttlvs := spcv1.TTLVS
	skr1, err := parseSKR1(ttlvs[Tag_SessionKey_R1])
	if err != nil {
		return err
	}

	appleD := d.AppleD{} //TODO: pass from parameter
	ask := []byte{}      //TODO:

	r2 := ttlvs[Tag_R2]
	dask, err := appleD.Compute(r2.Value, ask) //TODO: pass r2Block instead of r2.Block.value
	if err != nil {
		return err
	}

	DecryptedSKR1Payload, err := decryptSKR1Payload(*skr1, dask)
	if err != nil {
		return err
	}

	//Check the integrity of this SPC message
	skr1int := ttlvs[Tag_SessionKey_R1_integrity]

	if !reflect.DeepEqual(skr1int.Value, DecryptedSKR1Payload.IntegrityBytes) {
		return errors.New(" Check the integrity of the SPC failed.")
	}

	fmt.Printf("DASk Value:\n\t%s\n\n", hex.EncodeToString(dask))
	fmt.Printf("SPC SK Value:\n\t%s\n\n", hex.EncodeToString(DecryptedSKR1Payload.SK))
	fmt.Printf("SPC [SK..R1] IV Value:\n\t%s\n\n", hex.EncodeToString(skr1.IV))
	//fmt.Printf("SPC R1 Value:\n%s\n\n",hex.EncodeToString(DecryptedSKR1Payload.R1))

	return nil
}

func ParseSPCV1(playback []byte, pem []byte) (*SPCContainer, error) {
	spcContainer := ParseSPCContainer(playback)

	spck, err := decryptSPCK(pem, spcContainer.EncryptedAesKey)
	if err != nil {
		return nil, err
	}

	spcPayloadRow, err := decryptSPCpayload(spcContainer, spck)
	if err != nil {
		return nil, err
	}

	printDebugSPC(spcContainer)

	spcContainer.TTLVS = parseTLLVs(spcPayloadRow)

	return spcContainer, nil
}

func ParseSPCContainer(playback []byte) *SPCContainer {
	spcContainer := &SPCContainer{}
	spcContainer.Version = binary.BigEndian.Uint32(playback[0:4])
	spcContainer.Reserved = playback[4:8]
	spcContainer.AesKeyIV = playback[8:24]
	spcContainer.EncryptedAesKey = playback[24:152]
	spcContainer.CertificateHash = playback[152:172]
	spcContainer.SPCPlayloadLength = binary.BigEndian.Uint32(playback[172:176])
	spcContainer.SPCPlayload = playback[176 : 176+spcContainer.SPCPlayloadLength]

	return spcContainer
}

func parseTLLVs(spcpayload []byte) map[uint64]TLLVBlock {
	var m map[uint64]TLLVBlock
	m = make(map[uint64]TLLVBlock)

	for currentOffset := 0; currentOffset < len(spcpayload); {

		tag := binary.BigEndian.Uint64(spcpayload[currentOffset : currentOffset+Field_Tag_Length])
		currentOffset += Field_Tag_Length

		blockLength := binary.BigEndian.Uint32(spcpayload[currentOffset : currentOffset+Field_Block_Length])
		currentOffset += Field_Block_Length

		valueLength := binary.BigEndian.Uint32(spcpayload[currentOffset : currentOffset+Field_Value_Length])
		currentOffset += Field_Value_Length

		//paddingSize := blockLength - valueLength

		value := spcpayload[currentOffset : currentOffset+int(valueLength)]

		var skip bool
		switch tag {
		case Tag_SessionKey_R1:
			fmt.Printf("Tag_SessionKey_R1 -- %x\n", tag)
		case Tag_SessionKey_R1_integrity:
			fmt.Printf("Tag_SessionKey_R1_integrity -- %x\n", tag)
		case Tag_AntiReplaySeed:
			fmt.Printf("Tag_AntiReplaySeed -- %x\n", tag)
		case Tag_R2:
			fmt.Printf("Tag_R2 -- %x\n", tag)
		case Tag_ReturnRequest:
			fmt.Printf("Tag_ReturnRequest -- %x\n", tag)
		case Tag_AssetID:
			fmt.Printf("Tag_AssetID -- %x\n", tag)
		case Tag_TransactionID:
			fmt.Printf("Tag_TransactionID -- %x\n", tag)
		case Tag_ProtocolVersionsSupported:
			fmt.Printf("Tag_ProtocolVersionsSupported -- %x\n", tag)
		case Tag_ProtocolVersionUsed:
			fmt.Printf("Tag_ProtocolVersionUsed -- %x\n", tag)
		case Tag_treamingIndicator:
			fmt.Printf("Tag_treamingIndicator -- %x\n", tag)
		case Tag_kSKDServerClientReferenceTime:
			fmt.Printf("Tag_kSKDServerClientReferenceTime -- %x\n", tag)
		default:
			skip = true
		}

		if skip == false {
			fmt.Printf("Tag size:0x%x\n", valueLength)
			fmt.Printf("Tag length:0x%x\n", blockLength)
			fmt.Printf("Tag value:%s\n\n", hex.EncodeToString(value))

			tllvBlock := TLLVBlock{
				Tag:         tag,
				BlockLength: blockLength,
				ValueLength: valueLength,
				Value:       value,
			}

			m[tag] = tllvBlock

		}

		//TODO: paring ttlv
		currentOffset = currentOffset + int(blockLength)
	}

	return m
}

func parseSKR1(tllv TLLVBlock) (*SKR1TLLVBlock, error) {
	iv := tllv.Value[0:16]
	payload := tllv.Value[16:112]

	return &SKR1TLLVBlock{
		TLLVBlock: tllv,
		IV:        iv,
		Payload:   payload,
	}, nil
}

func decryptSKR1Payload(skr1 SKR1TLLVBlock, dask []byte) (*DecryptedSKR1Payload, error) {
	if skr1.Tag != Tag_SessionKey_R1 {
		return nil, errors.New("decryptSKR1 doesn't match Tag_SessionKey_R1 tag.")
	}

	decryptPayloadRow, err := aes.Decrypt(dask, skr1.IV, skr1.Payload)
	if err != nil {
		return nil, err
	}

	if len(decryptPayloadRow) != 96 {
		return nil, errors.New("Wrong decrypt payload size. Must be 96 bytes expected.")
	}

	d := &DecryptedSKR1Payload{
		SK:             decryptPayloadRow[0:16],
		HU:             decryptPayloadRow[16:36],
		R1:             decryptPayloadRow[36:80],
		IntegrityBytes: decryptPayloadRow[80:96],
	}

	return d, nil
}

func printDebugSPC(spcContainer *SPCContainer) {
	fmt.Println("========================= Begin SPC Data ===============================")
	fmt.Printf("SPC container size %+v\n", spcContainer.SPCPlayloadLength)

	fmt.Println("SPC Encryption Key -")
	fmt.Println(hex.EncodeToString(spcContainer.EncryptedAesKey))
	fmt.Println("SPC Encryption IV -")
	fmt.Println(hex.EncodeToString(spcContainer.AesKeyIV))
	fmt.Println("================ SPC TLLV List ================")

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
