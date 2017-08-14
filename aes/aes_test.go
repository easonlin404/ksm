package aes

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	var key = "1ae8ccd0e7985cc0b6203a55855a1034afc252980e970ca90e5202689f947ab9"
	var iv = "d58ce954203b7c9a9a9d467f59839249"

	keyByteAry, _ := hex.DecodeString(key)
	ivByteAry, _ := hex.DecodeString(iv)
	fmt.Println(len(keyByteAry))
	fmt.Println(len(ivByteAry))
	plainText := []byte("ABCDEFG")

	crypted, err := Encrypt(keyByteAry, ivByteAry, plainText)

	fmt.Println("crypted:%+v", crypted)
	enText := base64.StdEncoding.EncodeToString(crypted)
	assert.NoError(t, err)
	assert.Equal(t, "3iIEkNQUcSar6WP8QnW1Sg==", enText)
}

func TestDecrypt(t *testing.T) {
	var key = "1ae8ccd0e7985cc0b6203a55855a1034afc252980e970ca90e5202689f947ab9"
	var iv = "d58ce954203b7c9a9a9d467f59839249"

	keyByteAry, _ := hex.DecodeString(key)
	ivByteAry, _ := hex.DecodeString(iv)

	enBase64Str := "3iIEkNQUcSar6WP8QnW1Sg=="

	en, err := base64.StdEncoding.DecodeString(enBase64Str)
	if err != nil {
		fmt.Println(err)
	}

	plainText, err := Decrypt(keyByteAry, ivByteAry, en)

	assert.NoError(t, err)
	assert.Equal(t, "ABCDEFG", string(plainText))
}
