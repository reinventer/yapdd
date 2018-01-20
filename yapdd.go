// Package yapdd provides a client implementation of Yandex.Mail for Domain API
package yapdd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	clientTypeAdmin     = "admin"
	clientTypeRegistrar = "registrar"
)

type Client struct {
	httpCli    *http.Client
	clientType string
	pddToken   string
	oauthToken string
}

func New(token string, opts ...Option) *Client {
	cli := &Client{
		pddToken:   token,
		clientType: clientTypeAdmin,
	}

	for _, o := range opts {
		o(cli)
	}

	if cli.httpCli == nil {
		cli.httpCli = http.DefaultClient
	}

	return cli
}

type Option func(*Client)

func AsRegistrar(oauthToken string) Option {
	return func(cli *Client) {
		cli.clientType = clientTypeRegistrar
		cli.oauthToken = oauthToken
	}
}

func WithHTTPClient(httpCli *http.Client) Option {
	return func(cli *Client) {
		cli.httpCli = httpCli
	}
}

func (c *Client) getURL(section, action string, params *DNSRequestParams) string {
	u := fmt.Sprintf("https://pddimp.yandex.ru/api2/%s/%s/%s", c.clientType, section, action)
	if params != nil {
		u = u + "?" + url.Values(*params).Encode()
	}
	return u
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["PddToken"] = []string{c.pddToken}
	if c.clientType == clientTypeRegistrar {
		req.Header.Set("Authorization", "OAuth "+c.oauthToken)
	}

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(v)
}
