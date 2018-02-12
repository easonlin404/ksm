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

var pub = `-----BEGIN CERTIFICATE-----
MIIDfTCCAmWgAwIBAgIIboBT3GOPJ50wDQYJKoZIhvcNAQEFBQAwfTELMAkGA1UE
BhMCVVMxEzARBgNVBAoMCkFwcGxlIEluYy4xJjAkBgNVBAsMHUFwcGxlIENlcnRp
ZmljYXRpb24gQXV0aG9yaXR5MTEwLwYDVQQDDChEUk0gVGVjaG5vbG9naWVzIENl
cnRpZmljYXRpb24gQXV0aG9yaXR5MB4XDTExMTAxODAxNTcyMloXDTEzMTAxNzAx
NTcyMlowRjERMA8GA1UEAwwIUGFydG5lcjIxETAPBgNVBAsMCFBhcnRuZXIyMREw
DwYDVQQKDAhQYXJ0bmVyMjELMAkGA1UEBhMCVVMwgZ8wDQYJKoZIhvcNAQEBBQAD
gY0AMIGJAoGBALReAQ24va6MquxUkOyrVLE0vjc3rv3a16qndKKKGL6afpkN19xc
/cWw9A2W0FCSJYgkY+iyhGPAO4BLWe0QSonJz08GdeEMS2wmj87h8PLe6Yyu8Ida
3hH+snc7hv2bxX5AI72ETSQWlElky3tHLCYV2tqbTW4BGQZvvE4LfM+tAgMBAAGj
gbswgbgwJwYLKoZIhvdjZAYNAQMEGAGAgEeXuoURG4c6qSNQztlZmgq9dM3kTzAv
BgsqhkiG92NkBg0BBAQgAaWxaRPd6O3itrSL3iqhd3fcpUMMhDQTIebXMN1IfmQw
HQYDVR0OBBYEFDdUHOfoNQC1nqz9IzDvC/WJR1ssMAwGA1UdEwEB/wQCMAAwHwYD
VR0jBBgwFoAU6rShbWWjpF5JZST6HCRnrVoa0DMwDgYDVR0PAQH/BAQDAgUgMA0G
CSqGSIb3DQEBBQUAA4IBAQB4gFunl0sKeqGza5fdDd9Dj0O+rutFPqIFFLY60Qgl
jQdkzaHegMBqoON3I2KWRxgOeaewArmlgZjK8LoTv++HALB1Thf7N9AulyWVCg7J
i/hFKhTNpbNWBXSkKYn1QpcnohAnjLsrNED7R0b4A7z1yBhUjU96uRsKU+Dd6St9
XMlvvK49iSWNadfz7IictPrOjvHj4hRzepE43U5unevsth2FXu553LMCZw7gy4h9
IMYU4NZSWhf5z+wYpjtzYxdoqynjvihqFdGqYDC2drzpLLhaCXZhZUq2D1mXoQaY
6URsYkp6FRwIAx++KnIwE7Q3kK6s+5sRpKK4zZ0y0O9Z
-----END CERTIFICATE-----`

var pri = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQC0XgENuL2ujKrsVJDsq1SxNL43N6792teqp3Siihi+mn6ZDdfc
XP3FsPQNltBQkiWIJGPosoRjwDuAS1ntEEqJyc9PBnXhDEtsJo/O4fDy3umMrvCH
Wt4R/rJ3O4b9m8V+QCO9hE0kFpRJZMt7RywmFdram01uARkGb7xOC3zPrQIDAQAB
AoGBAIO+vkpFjNd4jEi/pHQa2WvuuJogpENsnGdclYc8E8L1mk81m1ys1/iUvk9G
v7Z6acu9uPR5oNYzzcJyR6cvZSFxtGIZnWNdDOAB71b+YqMvj3lr6MgUdMUgUfxZ
EDXLEhIoVzyQWIt+f6hjSG/hzyw+Jglo4ogCWPsV3S6UG2WBAkEA5HPddGIUa34k
2/EGQqyCAo4VYlCUdCFTp9+eFIUedequgsSIZhgblT+FSvMPYARuG/ywLoOivRy1
dFl0dIB1sQJBAModyMskK0r312kro+URq8VxlwwY0fv2rF1aS0/clQUw5OH/OxEn
Dgz3l3PNTXDCcQDh9wyEZV0SgIp7SYCDrL0CQEo8HEolVN1ZMEEIITCpPdX2tZws
8xCJg9WZJJUmbK+EgxCbLHeAffYRng6szOI2jlEp21ZCEC/DlHMqXl09IQECQGSn
EoC/oWOzKy4v0m3YL/+iwsL+dUwSGuJefhTmV7v/DmzRixvOpDum7WB5BDC8VERJ
Q5uTL1t7RFIydXcvm80CQH/E17mWT66PPeqloAfSH/5tJyak2gagkuFnMh779JRF
rl5YIIiAh+q5DkcjWw6eni5O4+UuwXRp29vZaxmDlIE=
-----END RSA PRIVATE KEY-----`
