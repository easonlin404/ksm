package ksm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTRLLVBlock(t *testing.T) {

	value := []byte{0x54, 0xa1, 0x6b, 0xe0, 0x13, 0x7e, 0xf2, 0x59, 0xab, 0x3e, 0x4f, 0xc7, 0x96, 0x90, 0x82, 0x5f}

	block := NewTLLVBlock(tagSessionKeyR1, value)

	assert.NotNil(t, block, "expect not nil")
	_, err := block.Serialize()

	assert.NoError(t, err)
}

func TestTLLVBlock_Serialize_error(t *testing.T) {
	b1 := NewTLLVBlock(tagSessionKeyR1, []byte{})
	_, err1 := b1.Serialize()
	assert.Error(t, err1)

	b2 := NewTLLVBlock(0, []byte{})
	_, err2 := b2.Serialize()
	assert.Error(t, err2)
}
