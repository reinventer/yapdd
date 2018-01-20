package yapdd

import (
	"context"
	"net/http"
	"strconv"
)

type DNSRecordType string

const (
	DNSTypeSRV   DNSRecordType = "SRV"
	DNSTypeTXT   DNSRecordType = "TXT"
	DNSTypeNS    DNSRecordType = "NS"
	DNSTypeMX    DNSRecordType = "MX"
	DNSTypeSOA   DNSRecordType = "SOA"
	DNSTypeA     DNSRecordType = "A"
	DNSTypeAAAA  DNSRecordType = "AAAA"
	DNSTypeCNAME DNSRecordType = "CNAME"
)

type DNSPriority struct {
	value uint16
	ok    bool
}

func (p *DNSPriority) UnmarshalJSON(b []byte) error {
	if string(b) == `""` {
		p.value = 0
		p.ok = false
		return nil
	}

	i, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	p.value = uint16(i)
	p.ok = true
	return nil
}

func (p *DNSPriority) Get() (uint16, bool) {
	return p.value, p.ok
}

type DNSRecord struct {
	ID        uint32        `json:"record_id"`
	Type      DNSRecordType `json:"type"`
	Domain    string        `json:"domain"`
	Subdomain string        `json:"subdomain"`
	FQDN      string        `json:"fqdn"`
	TTL       uint32        `json:"ttl"`
	Content   string        `json:"content"`
	Priority  DNSPriority   `json:"priority"`
	Operation string        `json:"operation"`
}

type DNSResponse struct {
	Domain   string     `json:"domain"`
	RecordID uint32     `json:"record_id"`
	Record   *DNSRecord `json:"record"`
	Success  string     `json:"success"`
	Error    string     `json:"error"`
}

type DNSListResponse struct {
	Domain  string       `json:"domain"`
	Records []*DNSRecord `json:"records"`
	Success string       `json:"success"`
	Error   string       `json:"error"`
}

func (c *Client) DNSAdd(ctx context.Context, domain string, recordType DNSRecordType, params *DNSRequestParams) (*DNSResponse, error) {
	params = params.recordType(recordType).domain(domain)
	req, err := http.NewRequest(
		http.MethodPost,
		c.getURL("dns", "add", nil),
		params.body(),
	)
	if err != nil {
		return nil, err
	}

	var r DNSResponse
	err = c.do(ctx, req, &r)
	return &r, err
}

func (c *Client) DNSList(ctx context.Context, domain string) (*DNSListResponse, error) {
	params := NewDNSParams().domain(domain)
	req, err := http.NewRequest(
		http.MethodGet,
		c.getURL("dns", "list", params),
		nil,
	)
	if err != nil {
		return nil, err
	}

	var r DNSListResponse
	err = c.do(ctx, req, &r)
	return &r, err
}

func (c *Client) DNSEdit(ctx context.Context, domain string, recordID uint32, params *DNSRequestParams) (*DNSResponse, error) {
	params = params.recordID(recordID).domain(domain)
	req, err := http.NewRequest(
		http.MethodPost,
		c.getURL("dns", "edit", nil),
		params.body(),
	)
	if err != nil {
		return nil, err
	}

	var r DNSResponse
	err = c.do(ctx, req, &r)
	return &r, err
}

func (c *Client) DNSDel(ctx context.Context, domain string, recordID uint32) (*DNSResponse, error) {
	params := NewDNSParams().recordID(recordID).domain(domain)
	req, err := http.NewRequest(
		http.MethodPost,
		c.getURL("dns", "del", nil),
		params.body(),
	)
	if err != nil {
		return nil, err
	}

	var r DNSResponse
	err = c.do(ctx, req, &r)
	return &r, err
}
