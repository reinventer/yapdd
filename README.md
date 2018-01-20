# Yapdd

[![Build Status](https://travis-ci.org/reinventer/yapdd.svg?branch=master)](https://travis-ci.org/reinventer/yapdd)

Yapdd is a golang package that provides a client implementation of [Yandex.Mail for Domain API](https://tech.yandex.com/domain/doc/about-docpage/).

Work still in progress.

## API supports
- [x] Managing DNS
- [ ] Managing DKIM
- [ ] Managing domain mailboxes
- [ ] Managing domain mailing lists
- [ ] Domain management
- [ ] Importing email
- [ ] Managing domain administrator proxies

## Example

```go
cli := yapdd.New("PddToken")
rec, err := cli.DNSAdd(
	context.Background(),
	"domain.com",
	yapdd.DNSTypeCNAME,
	yapdd.NewDNSParams().Subdomain("www").Content("domain.com"),
)
```

**Important note**: http.DefaultClient is used in package by default. Please replace the HTTP client if you want to use yapdd in production.
For example:
```go
httpCli:=&http.Client{Timeout: time.Second}
cli := yapdd.New("PddToken", yapdd.WithHTTPClient(httpCli))
```