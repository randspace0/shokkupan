package auth

import (
	"crypto/tls"
	"net/http"
	"testing"
)

func TestIsSecureRequest(t *testing.T) {
	plain := &http.Request{Header: http.Header{}}
	if isSecureRequest(plain) {
		t.Fatal("plain http must not be secure")
	}
	proxied := &http.Request{Header: http.Header{"X-Forwarded-Proto": {"https"}}}
	if !isSecureRequest(proxied) {
		t.Fatal("X-Forwarded-Proto https must be secure")
	}
	direct := &http.Request{Header: http.Header{}, TLS: &tls.ConnectionState{}}
	if !isSecureRequest(direct) {
		t.Fatal("direct TLS must be secure")
	}
}
