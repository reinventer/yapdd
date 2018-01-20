package yapdd

import (
	"net/http"
)

type httpTransportMock struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (m *httpTransportMock) RoundTrip(r *http.Request) (*http.Response, error) {
	m.request = r
	return m.response, m.err
}
