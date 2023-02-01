package httping

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type HttpResponse struct {
	Status          string
	Host            string
	ResponseHeaders map[string]string
	Latency         int64
}


// ParseURL will parse a URL from a string
func ParseURL(url string, useHTTP bool) string {
	if strings.HasPrefix(url, "http") || strings.HasPrefix(url,"https") {
		return url
	}

	if useHTTP {
		return "http://"+url
	}
	return "https://"+url
}

// MakeRequest performs an HTTP request
func MakeRequest(url, headers string) (*HttpResponse, error) {
	var result *HttpResponse

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	start := time.Now()
	response, err := client.Get(url)
	end := time.Since(start).Milliseconds()
	if err != nil {
		return &HttpResponse{}, err
	}
	defer response.Body.Close()

	h := make(map[string]string)
	responseHeaders := strings.Split(headers, ",")
	for _, header := range responseHeaders {
		h[header] = response.Header.Get(header)
	}

	result = &HttpResponse {
		Status: response.Status,
		Host: response.Header.Get("host"),
		ResponseHeaders: h,
		Latency: end,
	}
	return result, nil
}

func ParseMap(m *map[string]string) *string {
	var result string
	for k, v := range *m {
		// just making the output pretty
		if len(*m) > 1 {
			result += fmt.Sprintf(" {%s:%s} ", k, v)
		} else {
			result += fmt.Sprintf(" %s:%s ", k, v)
		}
	}
	return &result
}
