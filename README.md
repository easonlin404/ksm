# ksm

[![Build Status](https://travis-ci.org/easonlin404/ksm.svg)](https://travis-ci.org/easonlin404/ksm)
[![codecov](https://codecov.io/gh/easonlin404/ksm/branch/master/graph/badge.svg)](https://codecov.io/gh/easonlin404/ksm)
[![Go Report Card](https://goreportcard.com/badge/github.com/easonlin404/ksm)](https://goreportcard.com/report/github.com/easonlin404/ksm)
[![GoDoc](https://godoc.org/github.com/easonlin404/ksm?status.svg)](https://godoc.org/github.com/easonlin404/ksm)


Apple FairPlay Streaming ([FPS](https://developer.apple.com/streaming/fps/)) securely delivers keys to Apple mobile devices, Apple TV, and Safari on OS X, which will enable playback of encrypted video content. 

This project is FairPlay Streaming Key Security Module written in Go (Golang).

## Usage

### Start using it

Download and install it:

```bash
$ go get github.com/easonlin404/ksm
```

### Verify ckc
Perform verification utility `verify_ckc` to test KSM implementation.
```
testdata/verify_ckc -s testdata/FPS/spc1.bin -c testdata/FPS/ckc1.bin
```

### Simple example

See [example/basic.go](example/basic.go)

```go
package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/easonlin404/ksm"
	"github.com/easonlin404/ksm/d"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {

	r := gin.Default()

	type SpcMessage struct {
		Spc     string `json:"spc" binding:"required"`
		AssetId string `json:"assetId"`
	}

	r.POST("/fps/rest/getLicense", func(c *gin.Context) {
		var spcMessage SpcMessage
		if err := c.ShouldBindWith(&spcMessage, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		fmt.Printf("%v\n", spcMessage)

		playback, err := base64.StdEncoding.DecodeString(spcMessage.Spc)
		checkError(err)

		k := &ksm.Ksm{
			Pub:       pub,
			Pri:       pri,
			Rck:       RandomContentKey{},
			DFunction: d.AppleD{},
			Ask:       []byte{},
		}
		ckc, err2 := k.GenCKC(playback)
		checkError(err2)

		ckcBase64 := base64.StdEncoding.EncodeToString(ckc)
		c.JSON(http.StatusOK, gin.H{
			"ckc": ckcBase64,
		})
	})

	r.Run(":8080")

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type RandomContentKey struct {
}

func (RandomContentKey) FetchContentKey(assetId []byte) ([]byte, []byte, error) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	rand.Read(key)
	rand.Read(iv)
	return key, iv, nil
}

func (RandomContentKey) FetchContentKeyDuration(assetId []byte) (*ksm.CkcContentKeyDurationBlock, error) {

	LeaseDuration := rand.Uint32()  // The duration of the lease, if any, in seconds.
	RentalDuration := rand.Uint32() // The duration of the rental, if any, in seconds.

	return ksm.NewCkcContentKeyDurationBlock(LeaseDuration, RentalDuration), nil
}
