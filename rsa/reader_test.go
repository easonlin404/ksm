package rsa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileReader_Read(t *testing.T) {

	keyReader := FileReader{FileName: "../testdata/Development Credentials/dev_private_key.pem"}

	b, err := keyReader.ReadPem()
	assert.NoError(t, err)
	assert.Equal(t, dev_private_key, b)

}

var dev_private_key = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`)
