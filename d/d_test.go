package d

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppleD_Compute(t *testing.T) {

	r2 := []byte{}
	ask := []byte{}

	appleD := AppleD{}
	d, err := appleD.Compute(r2, ask)

	assert.NoError(t, err)
	fmt.Println(d)
}
