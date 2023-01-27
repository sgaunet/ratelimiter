package main

import (
	"net/http"
	"testing"
)

func TestGetIP(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://www.google.fr", nil)
	req.Header.Add("X-FORWARDED-FOR", "192.168.0.148, 10.20.130.02, 127.0.0.1")
	if GetIP(req) != "192.168.0.148" {
		t.Error(GetIP(req))
	}
	req, _ = http.NewRequest("GET", "http://www.google.fr", nil)
	if GetIP(req) != "" {
		t.Error(GetIP(req))
	}
}
