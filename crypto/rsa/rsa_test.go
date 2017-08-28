package rsa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)

var cert = []byte(`-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`)

func TestEncryptAndDecrypt(t *testing.T) {
	origData := []byte("plainText")

	en, err := EncryptByCert(cert, origData)
	assert.NoError(t, err)

	de, err := Decrypt(privateKey, en)
	assert.NoError(t, err)

	assert.Equal(t, origData, de)
}

func TestEncryptAndDecryptFromFile(t *testing.T) {
	origData := []byte("plainText")

	pemReader := FileReader{FileName: "../../testdata/Development Credentials/dev_private_key.pem"}
	derReader := FileReader{FileName: "../../testdata/Development Credentials/certificate.pem"}

	der, err := derReader.ReadPem()
	assert.NoError(t, err)

	pem, err := pemReader.ReadPem()
	assert.NoError(t, err)

	en, err := EncryptByCert(der, origData)
	assert.NoError(t, err)

	de, err := Decrypt(pem, en)
	assert.NoError(t, err)

	assert.Equal(t, origData, de)
}
