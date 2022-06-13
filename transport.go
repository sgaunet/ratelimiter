package main

import (
	"net/http"
)

type myTransport struct {
	// Uncomment this if you want to capture the transport
	// CapturedTransport http.RoundTripper
}

func (t *myTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(req)
	// or, if you captured the transport
	// response, err := t.CapturedTransport.RoundTrip(request)

	// The httputil package provides a DumpResponse() func that will copy the
	// contents of the body into a []byte and return it. It also wraps it in an
	// ioutil.NopCloser and sets up the response to be passed on to the client.
	// body, err := httputil.DumpResponse(response, true)
	// if err != nil {
	// 	// copying the response body did not work
	// 	return nil, err
	// }

	// You may want to check the Content-Type header to decide how to deal with
	// the body. In this case, we're assuming it's text.
	if response == nil {
		formatLog := "req.RemoteAddr=%s req.Host=%s req.URL.Path=%s req.URL.Query()=%s req.URL.Scheme=%s Error=%s\n"
		log.Errorf(formatLog, GetIP(req), req.Host, req.URL.Path, req.URL.Query(), req.URL.Scheme, err.Error())
	} else {
		formatLog := "req.RemoteAddr=%s req.Host=%s req.URL.Path=%s req.URL.Query()=%s StatusCode=%d\n"
		log.Infof(formatLog, GetIP(req), req.Host, req.URL.Path, req.URL.Query(), response.StatusCode)
	}

	return response, err
}
