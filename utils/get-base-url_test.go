package utils

import (
	"crypto/tls"
	"net/http"
	"testing"
)

func TestGetBaseURL(t *testing.T) {
	req := &http.Request{
		Host: "example.com",
		TLS:  nil,
	}
	got := GetBaseURL(req)
	want := "http://example.com"
	if got != want {
		t.Errorf("GetBaseURL() = %q, want %q", got, want)
	}

	reqTLS := &http.Request{
		Host: "secure.com",
		TLS:  &tls.ConnectionState{},
	}
	gotTLS := GetBaseURL(reqTLS)
	wantTLS := "https://secure.com"
	if gotTLS != wantTLS {
		t.Errorf("GetBaseURL() = %q, want %q", gotTLS, wantTLS)
	}
}
