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

type HttpStatistics struct {
	Count200 			 int
	Count300 			 int
	Count400 			 int
	Count500 			 int
	Other 				 int
	AverageLatency int64
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

func ParseHeader(m *map[string]string) *string {
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

func CalculateStatistics(responses []*HttpResponse) *HttpStatistics {
	var stats HttpStatistics
	var totalLatency int64

	for _, l := range responses {
		totalLatency += l.Latency

		switch strings.Trim(l.Status,"\n") {
		case "200 OK":
			stats.Count200++
		case "300":
			stats.Count300++
		case "400":
			stats.Count400++
		case "500":
			stats.Count500++
		default:
			stats.Other++
		}
	}
	stats.AverageLatency = totalLatency / int64(len(responses))

	return &stats
}