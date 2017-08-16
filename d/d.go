package d

import (
	"encoding/hex"
	"errors"
)

//DerivedApplicationSecretKey
type D interface {
	//
	Compute(R2 []byte, ask []byte) ([]byte, error)
}

type AppleD struct {
}

func (d AppleD) Compute(R2 []byte, ask []byte) ([]byte, error) {

	if len(R2) == 0 {
		errors.New("R2 block doesn't exist.")
	}

	b, err := hex.DecodeString("d87ce7a26081de2e8eb8acef3a6dc179")

	if len(b) != 16 {
		errors.New("ask key length doesn't equal 16")
	}

	return b, err

}
