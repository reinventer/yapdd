package yapdd

import (
	"io"
	"net/url"
	"strconv"
	"strings"
)

type DNSRequestParams url.Values

func NewDNSParams() *DNSRequestParams {
	p := DNSRequestParams(url.Values{})
	return &p
}

func (p *DNSRequestParams) domain(domain string) *DNSRequestParams {
	url.Values(*p).Set("domain", domain)
	return p
}

func (p *DNSRequestParams) recordType(t DNSRecordType) *DNSRequestParams {
	url.Values(*p).Set("type", string(t))
	return p
}

func (p *DNSRequestParams) recordID(id uint32) *DNSRequestParams {
	url.Values(*p).Set("record_id", strconv.Itoa(int(id)))
	return p
}

func (p *DNSRequestParams) AdminMail(email string) *DNSRequestParams {
	url.Values(*p).Set("admin_mail", email)
	return p
}

func (p *DNSRequestParams) Content(content string) *DNSRequestParams {
	url.Values(*p).Set("content", content)
	return p
}

func (p *DNSRequestParams) Priority(priority uint16) *DNSRequestParams {
	url.Values(*p).Set("priority", strconv.Itoa(int(priority)))
	return p
}

func (p *DNSRequestParams) Weight(weight uint16) *DNSRequestParams {
	url.Values(*p).Set("weight", strconv.Itoa(int(weight)))
	return p
}

func (p *DNSRequestParams) Port(port uint16) *DNSRequestParams {
	url.Values(*p).Set("port", strconv.Itoa(int(port)))
	return p
}

func (p *DNSRequestParams) Target(target string) *DNSRequestParams {
	url.Values(*p).Set("target", target)
	return p
}

func (p *DNSRequestParams) Subdomain(subdomain string) *DNSRequestParams {
	url.Values(*p).Set("subdomain", subdomain)
	return p
}

func (p *DNSRequestParams) TTL(ttl uint32) *DNSRequestParams {
	url.Values(*p).Set("ttl", strconv.Itoa(int(ttl)))
	return p
}

func (p *DNSRequestParams) Refresh(refresh uint32) *DNSRequestParams {
	url.Values(*p).Set("refresh", strconv.Itoa(int(refresh)))
	return p
}

func (p *DNSRequestParams) SetRetry(retry uint32) *DNSRequestParams {
	url.Values(*p).Set("retry", strconv.Itoa(int(retry)))
	return p
}

func (p *DNSRequestParams) Expire(expire uint16) *DNSRequestParams {
	url.Values(*p).Set("expire", strconv.Itoa(int(expire)))
	return p
}

func (p *DNSRequestParams) NegCache(negcache uint32) *DNSRequestParams {
	url.Values(*p).Set("neg_cache", strconv.Itoa(int(negcache)))
	return p
}

func (p *DNSRequestParams) body() io.Reader {
	return strings.NewReader(url.Values(*p).Encode())
}
