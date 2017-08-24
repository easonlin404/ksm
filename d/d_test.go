package d

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppleD_Compute(t *testing.T) {
	r2 := []byte{}
	ask := []byte{}

	appleD := AppleD{}
	_, err := appleD.Compute(r2, ask)

	assert.NoError(t, err)
}
