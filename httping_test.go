package httping_test

import (
	"testing"

	"github.com/rnemeth90/httping"
)

func TestParseURL(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		insecure bool
		expect   string
	}{
		{name: "insecure", url: "http://www.google.com", insecure: true, expect: "http://www.google.com"},
		{name: "secure", url: "https://www.google.com", insecure: false, expect: "https://www.google.com"},
		{name: "no-protocol-insecure", url: "www.google.com", insecure: true, expect: "http://www.google.com"},
		{name: "no-protocol-secure", url: "www.google.com", insecure: false, expect: "https://www.google.com"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got := httping.ParseURL(test.url, test.insecure)

			if got != test.expect {
				t.Errorf("expected %s\n, got %s\n", test.expect, got)
			}
		})
	}
}

func TestParseHeader(t *testing.T) {
	testCases := []struct {
		name    string
		headers map[string]string
		expect  string
	}{
		{name: "ManyHeaders",
			headers: map[string]string{
				"host":         "server01",
				"Content-Type": "application/json",
			},
			expect: " {host:server01}  {Content-Type:application/json} ",
		},
		{name: "OneHeader",
			headers: map[string]string{
				"{host": "server01}",
			},
			expect: " {host:server01} ",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got := httping.ParseHeader(&test.headers)

			if test.expect != *got {
				t.Errorf("expected %s\n, got %s\n", test.expect, *got)
			}
		})
	}
}

func TestCalculateStatistics(t *testing.T) {
	testCases := []struct{
		name string
		responses []httping.HttpResponse
		expect httping.HTTPStatistics
	}{
		{"Test200s",responses: []httping.HttpResponse{}, expect: httping.HTTPStatistics{}},
	}
}