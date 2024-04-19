package httping_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
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
				"Content-Type": "application/json",
				"host":         "server01",
			},
			expect: " {Content-Type:application/json}  {host:server01} ",
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

			if test.expect != got {
				t.Errorf("expected %s\n, got %s\n", test.expect, got)
			}
		})
	}
}

func TestCalculateStatistics(t *testing.T) {
	testCases := []struct {
		name      string
		responses []*httping.HttpResponse
		expect    *httping.HTTPStatistics
	}{
		{name: "Test200s",
			responses: []*httping.HttpResponse{
				{
					Status:          200,
					Host:            "localhost",
					ResponseHeaders: nil,
					Latency:         123,
				},
				{
					Status:          302,
					Host:            "localhost",
					ResponseHeaders: nil,
					Latency:         100,
				},
			},
			expect: &httping.HTTPStatistics{
				Count200:       1,
				Count204:       0,
				Count201:       0,
				Count301:       0,
				Count302:       1,
				Count304:       0,
				Count400:       0,
				Count401:       0,
				Count403:       0,
				Count404:       0,
				Count500:       0,
				Count502:       0,
				Count503:       0,
				Count504:       0,
				Other:          0,
				AverageLatency: 111,
				MaxLatency:     123,
				MinLatency:     100,
			}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got := httping.CalculateStatistics(test.responses)

			if !reflect.DeepEqual(test.expect, got) {
				t.Errorf("expected %q\n, got %q\n", test.expect, got)
			}
		})
	}
}

func TestMakeRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("host", "tester")
		w.WriteHeader(http.StatusTeapot)
	}))

	client := http.Client{}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	expected := httping.HttpResponse{
		Status:          resp.StatusCode,
		Host:            "tester",
		ResponseHeaders: nil,
		Latency:         0,
	}

	got, err := httping.MakeRequest(false, "", server.URL, "host")
	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(got, expected) {
		t.Errorf("expected response %q\n, but got %q\n", expected.Host, got)
	}
}
