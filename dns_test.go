package yapdd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestClient_DNSAdd(t *testing.T) {
	cases := []struct {
		name           string
		asRegistrar    bool
		httpResponse   *http.Response
		httpErr        error
		expErr         error
		expResponse    *DNSResponse
		expHTTPRequest *http.Request
	}{
		{
			name: "success",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`
					{
					  "domain": "domain.com", 
					  "record":
					  { 
						"record_id": 1,
						"type": "A",
						"domain": "domain.com",
						"subdomain": "www",
						"fqdn": "www.domain.com",
						"content": "1.2.3.4",
						"ttl": 900,
						"priority": ""
					  },
					  "success": "ok"
					}
				`)),
			},
			expResponse: &DNSResponse{
				Domain: "domain.com",
				Record: &DNSRecord{
					ID:        1,
					Type:      DNSTypeA,
					Domain:    "domain.com",
					Subdomain: "www",
					FQDN:      "www.domain.com",
					Content:   "1.2.3.4",
					TTL:       900,
					Priority:  DNSPriority{},
				},
				Success: "ok",
			},
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/add",
				"content=1.2.3.4&domain=domain.com&subdomain=www&type=A",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name:        "success as registrar",
			asRegistrar: true,
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`
					{
					  "domain": "domain.com", 
					  "record":
					  { 
						"record_id": 1,
						"type": "A",
						"domain": "domain.com",
						"subdomain": "www",
						"fqdn": "www.domain.com",
						"content": "1.2.3.4",
						"ttl": 900,
						"priority": ""
					  },
					  "success": "ok"
					}
				`)),
			},
			expResponse: &DNSResponse{
				Domain: "domain.com",
				Record: &DNSRecord{
					ID:        1,
					Type:      DNSTypeA,
					Domain:    "domain.com",
					Subdomain: "www",
					FQDN:      "www.domain.com",
					Content:   "1.2.3.4",
					TTL:       900,
					Priority:  DNSPriority{},
				},
				Success: "ok",
			},
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/registrar/dns/add",
				"content=1.2.3.4&domain=domain.com&subdomain=www&type=A",
				map[string][]string{
					"PddToken":      {"token"},
					"Authorization": {"OAuth oauthToken"},
					"Content-Type":  {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name:        "fail: transport returned error",
			httpErr:     errors.New("fail"),
			expResponse: &DNSResponse{},
			expErr:      errors.New("Post https://pddimp.yandex.ru/api2/admin/dns/add: fail"),
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/add",
				"content=1.2.3.4&domain=domain.com&subdomain=www&type=A",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name: "fail: bad json in response",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("bad json")),
			},
			expResponse: &DNSResponse{},
			expErr:      errors.New("invalid character 'b' looking for beginning of value"),
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/add",
				"content=1.2.3.4&domain=domain.com&subdomain=www&type=A",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(*testing.T) {
			transport := &httpTransportMock{
				response: tc.httpResponse,
				err:      tc.httpErr,
			}
			httpClient := &http.Client{
				Transport: transport,
			}

			p := []Option{
				WithHTTPClient(httpClient),
			}

			if tc.asRegistrar {
				p = append(p, AsRegistrar("oauthToken"))
			}

			cli := New("token", p...)

			response, err := cli.DNSAdd(
				context.Background(),
				"domain.com",
				DNSTypeA,
				NewDNSParams().Subdomain("www").Content("1.2.3.4"),
			)
			if fmt.Sprint(tc.expErr) != fmt.Sprint(err) {
				t.Errorf("expected error: %v, got: %v", tc.expErr, err)
			}
			if !reflect.DeepEqual(tc.expResponse, response) {
				t.Errorf("expected response: %+v, got: %+v", tc.expResponse, response)
			}

			ok, err := requestsEqual(tc.expHTTPRequest, transport.request)
			if err != nil {
				t.Fatalf("error reading body of request: %s", err)
			}
			if !ok {
				t.Errorf("expected request:\n%+v,\ngot:\n%+v", tc.expHTTPRequest, transport.request)
			}
		})
	}
}

func TestClient_DNSList(t *testing.T) {
	cases := []struct {
		name           string
		httpResponse   *http.Response
		httpErr        error
		expErr         error
		expResponse    *DNSListResponse
		expHTTPRequest *http.Request
	}{
		{
			name: "success",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`
					{
					  "domain": "domain.com", 
					  "records": [
						  { 
							"record_id": 1,
							"type": "A",
							"domain": "domain.com",
							"subdomain": "www",
							"fqdn": "www.domain.com",
							"content": "1.2.3.4",
							"ttl": 900,
							"priority": ""
						  },
						  { 
							"record_id": 2,
							"type": "SRV",
							"domain": "domain.com",
							"subdomain": "_xmpp-server._tcp",
							"content": "domain-xmpp.domain.com",
							"priority": 20
						  }
					  ],
					  "success": "ok"
					}
				`)),
			},
			expResponse: &DNSListResponse{
				Domain: "domain.com",
				Records: []*DNSRecord{
					{
						ID:        1,
						Type:      DNSTypeA,
						Domain:    "domain.com",
						Subdomain: "www",
						FQDN:      "www.domain.com",
						Content:   "1.2.3.4",
						TTL:       900,
						Priority:  DNSPriority{},
					},
					{
						ID:        2,
						Type:      DNSTypeSRV,
						Domain:    "domain.com",
						Subdomain: "_xmpp-server._tcp",
						Content:   "domain-xmpp.domain.com",
						Priority:  DNSPriority{value: 20, ok: true},
					},
				},
				Success: "ok",
			},
			expHTTPRequest: getRequest(
				t,
				http.MethodGet,
				"https://pddimp.yandex.ru/api2/admin/dns/list?domain=domain.com",
				"",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name:        "fail: transport returned error",
			httpErr:     errors.New("fail"),
			expResponse: &DNSListResponse{},
			expErr:      errors.New("Get https://pddimp.yandex.ru/api2/admin/dns/list?domain=domain.com: fail"),
			expHTTPRequest: getRequest(
				t,
				http.MethodGet,
				"https://pddimp.yandex.ru/api2/admin/dns/list?domain=domain.com",
				"",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name: "fail: bad json in response",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("bad json")),
			},
			expResponse: &DNSListResponse{},
			expErr:      errors.New("invalid character 'b' looking for beginning of value"),
			expHTTPRequest: getRequest(
				t,
				http.MethodGet,
				"https://pddimp.yandex.ru/api2/admin/dns/list?domain=domain.com",
				"",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(*testing.T) {
			transport := &httpTransportMock{
				response: tc.httpResponse,
				err:      tc.httpErr,
			}
			httpClient := &http.Client{
				Transport: transport,
			}

			cli := New("token", WithHTTPClient(httpClient))

			response, err := cli.DNSList(
				context.Background(),
				"domain.com",
			)
			if fmt.Sprint(tc.expErr) != fmt.Sprint(err) {
				t.Errorf("expected error: %v, got: %v", tc.expErr, err)
			}
			if !reflect.DeepEqual(tc.expResponse, response) {
				t.Errorf("expected response: %+v, got: %+v", tc.expResponse, response)
			}

			ok, err := requestsEqual(tc.expHTTPRequest, transport.request)
			if err != nil {
				t.Fatalf("error reading body of request: %s", err)
			}
			if !ok {
				t.Errorf("expected request:\n%+v,\ngot:\n%+v", tc.expHTTPRequest, transport.request)
			}
		})
	}
}

func TestClient_DNSEdit(t *testing.T) {
	cases := []struct {
		name           string
		httpResponse   *http.Response
		httpErr        error
		expErr         error
		expResponse    *DNSResponse
		expHTTPRequest *http.Request
	}{
		{
			name: "success",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`
					{
					  "domain": "domain.com", 
					  "record":
					  { 
						"record_id": 1,
						"type": "A",
						"domain": "domain.com",
						"subdomain": "www",
						"fqdn": "www.domain.com",
						"content": "1.2.3.4",
						"ttl": 900,
						"priority": "",
						"operation": "editing"
					  },
					  "success": "ok"
					}
				`)),
			},
			expResponse: &DNSResponse{
				Domain: "domain.com",
				Record: &DNSRecord{
					ID:        1,
					Type:      DNSTypeA,
					Domain:    "domain.com",
					Subdomain: "www",
					FQDN:      "www.domain.com",
					Content:   "1.2.3.4",
					TTL:       900,
					Priority:  DNSPriority{},
					Operation: "editing",
				},
				Success: "ok",
			},
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/edit",
				"content=1.2.3.4&domain=domain.com&record_id=1&subdomain=www",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name:        "fail: transport returned error",
			httpErr:     errors.New("fail"),
			expResponse: &DNSResponse{},
			expErr:      errors.New("Post https://pddimp.yandex.ru/api2/admin/dns/edit: fail"),
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/edit",
				"content=1.2.3.4&domain=domain.com&record_id=1&subdomain=www",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name: "fail: unexpected http status",
			httpResponse: &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Status:     "503 Service Unavailable",
				Body:       ioutil.NopCloser(strings.NewReader("")),
			},
			expResponse: &DNSResponse{},
			expErr:      errors.New("unexpected response status: 503 Service Unavailable"),
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/edit",
				"content=1.2.3.4&domain=domain.com&record_id=1&subdomain=www",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name: "fail: bad json in response",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("bad json")),
			},
			expResponse: &DNSResponse{},
			expErr:      errors.New("invalid character 'b' looking for beginning of value"),
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/edit",
				"content=1.2.3.4&domain=domain.com&record_id=1&subdomain=www",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(*testing.T) {
			transport := &httpTransportMock{
				response: tc.httpResponse,
				err:      tc.httpErr,
			}
			httpClient := &http.Client{
				Transport: transport,
			}

			cli := New("token", WithHTTPClient(httpClient))

			response, err := cli.DNSEdit(
				context.Background(),
				"domain.com",
				1,
				NewDNSParams().Subdomain("www").Content("1.2.3.4"),
			)
			if fmt.Sprint(tc.expErr) != fmt.Sprint(err) {
				t.Errorf("expected error: %v, got: %v", tc.expErr, err)
			}
			if !reflect.DeepEqual(tc.expResponse, response) {
				t.Errorf("expected response: %+v, got: %+v", tc.expResponse, response)
			}

			ok, err := requestsEqual(tc.expHTTPRequest, transport.request)
			if err != nil {
				t.Fatalf("error reading body of request: %s", err)
			}
			if !ok {
				t.Errorf("expected request:\n%+v,\ngot:\n%+v", tc.expHTTPRequest, transport.request)
			}
		})
	}
}

func TestClient_DNSDel(t *testing.T) {
	cases := []struct {
		name           string
		httpResponse   *http.Response
		httpErr        error
		expErr         error
		expResponse    *DNSResponse
		expHTTPRequest *http.Request
	}{
		{
			name: "success",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`
					{
					  "domain": "domain.com",
					  "record_id": 1,
					  "success": "ok"
					}
				`)),
			},
			expResponse: &DNSResponse{
				Domain:   "domain.com",
				RecordID: 1,
				Success:  "ok",
			},
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/del",
				"domain=domain.com&record_id=1",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name:        "fail: transport returned error",
			httpErr:     errors.New("fail"),
			expResponse: &DNSResponse{},
			expErr:      errors.New("Post https://pddimp.yandex.ru/api2/admin/dns/del: fail"),
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/del",
				"domain=domain.com&record_id=1",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
		{
			name: "fail: bad json in response",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("bad json")),
			},
			expResponse: &DNSResponse{},
			expErr:      errors.New("invalid character 'b' looking for beginning of value"),
			expHTTPRequest: getRequest(
				t,
				http.MethodPost,
				"https://pddimp.yandex.ru/api2/admin/dns/del",
				"domain=domain.com&record_id=1",
				map[string][]string{
					"PddToken":     {"token"},
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(*testing.T) {
			transport := &httpTransportMock{
				response: tc.httpResponse,
				err:      tc.httpErr,
			}
			httpClient := &http.Client{
				Transport: transport,
			}

			cli := New("token", WithHTTPClient(httpClient))

			response, err := cli.DNSDel(
				context.Background(),
				"domain.com",
				1,
			)
			if fmt.Sprint(tc.expErr) != fmt.Sprint(err) {
				t.Errorf("expected error: %v, got: %v", tc.expErr, err)
			}
			if !reflect.DeepEqual(tc.expResponse, response) {
				t.Errorf("expected response: %+v, got: %+v", tc.expResponse, response)
			}

			ok, err := requestsEqual(tc.expHTTPRequest, transport.request)
			if err != nil {
				t.Fatalf("error reading body of request: %s", err)
			}
			if !ok {
				t.Errorf("expected request:\n%+v,\ngot:\n%+v", tc.expHTTPRequest, transport.request)
			}
		})
	}
}

func getRequest(t *testing.T, method, url string, body string, headers http.Header) *http.Request {
	r, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("can't get request %s %s", method, url)
	}
	r.Header = headers
	return r.WithContext(context.Background())
}

// requestsEqual compares methods, URLs, headers and bodies of two requests
func requestsEqual(r1, r2 *http.Request) (bool, error) {
	if r1.Method == r2.Method &&
		r1.URL.String() == r2.URL.String() &&
		reflect.DeepEqual(r1.Header, r2.Header) {

		var (
			body1, body2 []byte
			err          error
		)

		if r1.Body != nil {
			body1, err = ioutil.ReadAll(r1.Body)
			if err != nil {
				return false, err
			}
		}

		if r2.Body != nil {
			body2, err = ioutil.ReadAll(r2.Body)
			if err != nil {
				return false, err
			}
		}

		if bytes.Equal(body1, body2) {
			return true, nil
		}
	}

	return false, nil
}
