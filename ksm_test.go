package ksm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenCKC(t *testing.T) {
	spcMessage:="00000001"

	err:=GenCKC(spcMessage)

	assert.NoError(t,err)
}
