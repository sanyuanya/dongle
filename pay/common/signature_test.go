package common

import "testing"

func TestExtractSignature(t *testing.T) {
	authorization := "WECHATPAY2-SHA256-RSA2048 mchid=\"3424\",nonce_str=\"234234\",timestamp=\"324234\",serial_no=\"345235\",signature=\"sdfsahfagsdhfgahsdfsdgfhsad\""
	signature, err := ExtractSignature(authorization)
	if err != nil {
		t.Errorf("TestExtractSignature failed, err: %v", err)
	}
	t.Logf("signature: %v", signature)

}
