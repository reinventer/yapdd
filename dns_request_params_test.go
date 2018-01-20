package yapdd

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestDNSRequestParams_SetAllParams(t *testing.T) {
	params := NewDNSParams().
		domain("domain.com").
		recordType(DNSTypeSRV).
		recordID(1).
		AdminMail("admin@domain.com").
		Content("content will be ignored").
		Content("content"). // last parameter will be used
		Expire(2).
		NegCache(3).
		Port(4).
		Priority(5).
		Refresh(6).
		SetRetry(7).
		Weight(8).
		Target("target").
		Subdomain("www").
		TTL(9)

	bodyReader := params.body()

	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		t.Fatalf("unexpected error while reading body: %s", err)
	}

	expBody := []byte("admin_mail=admin%40domain.com&content=content&domain=domain.com&expire=2&neg_cache=3&port=4&priority=5&record_id=1&refresh=6&retry=7&subdomain=www&target=target&ttl=9&type=SRV&weight=8")
	if !bytes.Equal(expBody, body) {
		t.Errorf("\nexpected body:\n%s\ngot:\n%s", expBody, body)
	}
}
