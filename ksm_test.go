package ksm

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"github.com/easonlin404/ksm/rsa"
)

type spcContainerTest struct {
	filePath string
	size     int

	SPCContainer
}

var spcContainerTests = []spcContainerTest{
	{"testdata/spc1.bin", 2688,
		SPCContainer{
			AesKeyIV: []byte{93, 22, 68, 234, 236, 17, 249, 131, 20, 117, 65, 228, 110, 235, 39, 116},
		}},
}

func TestGenCKC(t *testing.T) {
	for _, test := range spcContainerTests {
		f, err := os.Open(test.filePath)
		assert.NoError(t, err)
		defer f.Close()

		//spcMessage, err := ioutil.ReadAll(f)
		//assert.NoError(t, err)

	}

}

func TestParseSPCContainer(t *testing.T) {
	FileReader:=rsa.FileReader{FileName:"testdata/Development Credentials/dev_private_key.pem"}

	for _, test := range spcContainerTests {
		f, err := os.Open(test.filePath)
		defer f.Close()
		assert.NoError(t, err)


		spcMessage, err := ioutil.ReadAll(f)
		assert.NoError(t, err)

		spcContainer, err := ParseSPCContainer(spcMessage,&FileReader)
		assert.NoError(t, err)

		assert.Equal(t, test.AesKeyIV, spcContainer.AesKeyIV)

		fmt.Println(hex.EncodeToString(spcContainer.EncryptedAesKey))

	}
}
