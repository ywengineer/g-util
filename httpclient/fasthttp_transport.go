package httpclient

import (
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var fastClient = fasthttp.Client{
	Name:         "CustomFasthttpClient",
	ReadTimeout:  2 * time.Second,
	WriteTimeout: 1 * time.Second,
}

// Transport implements the estransport interface with
// the github.com/valyala/fasthttp HTTP client.
//
type fastHttpTransport struct {
}

func NewFastHttpTransport() *fastHttpTransport {
	return &fastHttpTransport{}
}

// RoundTrip performs the request and returns a response or error
//
func (t *fastHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	freq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(freq)

	fres := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(fres)

	t.copyRequest(freq, req)

	err := fastClient.Do(freq, fres)
	if err != nil {
		return nil, err
	}

	res := &http.Response{Header: make(http.Header)}
	t.copyResponse(res, fres)

	return res, nil
}

// copyRequest converts a http.Request to fasthttp.Request
//
func (t *fastHttpTransport) copyRequest(dst *fasthttp.Request, src *http.Request) *fasthttp.Request {
	if src.Method == "GET" && src.Body != nil {
		src.Method = "POST"
	}

	dst.SetHost(src.Host)
	dst.SetRequestURI(src.URL.String())

	dst.Header.SetRequestURI(src.URL.String())
	dst.Header.SetMethod(src.Method)

	for k, vv := range src.Header {
		for _, v := range vv {
			dst.Header.Set(k, v)
		}
	}

	if src.Body != nil {
		dst.SetBodyStream(src.Body, -1)
	}

	return dst
}

// copyResponse converts a http.Response to fasthttp.Response
//
func (t *fastHttpTransport) copyResponse(dst *http.Response, src *fasthttp.Response) *http.Response {
	dst.StatusCode = src.StatusCode()

	src.Header.VisitAll(func(k, v []byte) {
		dst.Header.Set(string(k), string(v))
	})

	// Cast to a string to make a copy seeing as src.Body() won't
	// be valid after the response is released back to the pool (fasthttp.ReleaseResponse).
	dst.Body = ioutil.NopCloser(strings.NewReader(string(src.Body())))

	return dst
}
