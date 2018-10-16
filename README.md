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

See [example/basic.go](example/basic.go). You have to implement your contentKey and D func if you want to use Apple FairPlay DRM.

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
			Rck:       RandomContentKey{}, // Use random content key for testing
			DFunction: d.AppleD{},         // Use D function provided by Apple Inc.
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

// Random content key
type RandomContentKey struct {
}

// Implement FetchContentKey func
func (RandomContentKey) FetchContentKey(assetId []byte) ([]byte, []byte, error) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	rand.Read(key)
	rand.Read(iv)
	return key, iv, nil
}

// Implement FetchContentKeyDuration func
func (RandomContentKey) FetchContentKeyDuration(assetId []byte) (*ksm.CkcContentKeyDurationBlock, error) {

	LeaseDuration := rand.Uint32()  // The duration of the lease, if any, in seconds.
	RentalDuration := rand.Uint32() // The duration of the rental, if any, in seconds.

	return ksm.NewCkcContentKeyDurationBlock(LeaseDuration, RentalDuration), nil
}
```

## FAQ
### How to send sample SPC data?

Our [example/basic.go](https://github.com/easonlin404/ksm/blob/master/example/basic.go) router has to accepted SPC parameter that is base64 encoding. And [testdata/spc1.bin](https://github.com/easonlin404/ksm/blob/master/testdata/FPS/spc1.bin) is binary file, so you have to encode to base64  string.


After run as local server, you can test sample POST request using [httpie](https://httpie.org/) tool(spc is spc1.bin base64 encoding):


```
http POST http://localhost:8080/fps/rest/getLicense spc=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
```

Will response:
```
HTTP/1.1 200 OK
Content-Length: 902
Content-Type: application/json; charset=utf-8
Date: Thu, 08 Mar 2018 05:51:02 GMT

{
    "ckc": "AAAAAQAAAAAqQk6qsqGlLq37TNsWDOnKAAACgKtZF0mIK6FPslvHP3oQdTZgepM5j4nnI/d//Rj+T4eTBfmKA4aVPwUr4bOozAhPvpL9QhXKIsGo9sEmtz0j0okdbnTTWCYAfLTHGwYviWMsR0J/vFtlFnxJF2yyLhDN6wu4/Be8Q7+B2duXsW7k5r6TF/gaKYbgTChM/3b+znqizRy06gNJr4d5HUKfXZ9FAmBtcsTUGbdMWjvKYJV8zPiYoayTkarXN6bDGMY2ToDyWBXIethDf4kBgdpOcR+qvXLntbI9daaEbb2XEU/nw4v6Dks6/sMMizb0rtY3fA4dCxEYeVxlAOLlLGLxswtRkZmNj1YQD/2HUGmuLBo7Q+17MGvRP06d+LG3YaDWbLeRPH5CVCwaxvmj1mFSLvMZffLgPSYZUp9fXaEeUnh8aD/mvaaFJavLmB7+uk+YggX+9F1HtjpFpjelvX5InP2L3mSYOl/eNdjzJbFw7l9BnR5KWL1wWiOmJHEWBw3Y+Tx5bLs/1QMasyTG+wAGxiXN+XyQW8frVO7e1xeBnuJplSEx+DK4Z9ZtveQG2l56OMsyw3KxUPTjCzHF5n1z6CzCSw5R57gA23mQ3nScPkDCbzQyaquZfCO80f6JW9hMGVEIFPnZVtUR1IemNKXK3fx+fZ9GdIN0EaNCp9fzB5GDC2dxrs9zhux3SvR5nKbSqHCX2b0S2S7Vn0LWZ7U4MtXMB2R8rsza03JEswAZ7iGKgl1J5L7TaADKz5cZvBkTR5P6QgHiesoZ9HMxQHtCyU4pjjGlFU081UKmsxefuuddPjwK93M/Xhf4wjOjcSSG5a1+Gg9FTi4tETdnMXVl/PuWjX1WFYCIJ1ZyKc2gZD2iEbw="
}
```


### How to verifying Key Security Module (KSM) Implementation?
[https://developer.apple.com/library/archive/technotes/tn2454/_index.html#//apple_ref/doc/uid/DTS40017630-CH1-VERIFYING_KEY_SECURITY_MODULE__KSM__IMPLEMENTATION](https://developer.apple.com/library/archive/technotes/tn2454/_index.html#//apple_ref/doc/uid/DTS40017630-CH1-VERIFYING_KEY_SECURITY_MODULE__KSM__IMPLEMENTATION)

### How to Debugging FairPlay Streaming?
[https://developer.apple.com/library/archive/technotes/tn2454/_index.html](https://developer.apple.com/library/archive/technotes/tn2454/_index.html)