// Package httping provides functionality to 'ping' web servers
// and gather statistics about the responses.
package httping

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HttpResponse represents the response received from an HTTP request.
// It includes status code, headers, latency, and any error encountered.
type HttpResponse struct {
	Status          int
	Host            string
	ResponseHeaders map[string]string
	Latency         int64
	Error           string
}

// HTTPStatistics aggregates statistics about a series of HTTP responses,
// including counts of different response codes and latency metrics.
type HTTPStatistics struct {
	Count200       int
	Count201       int
	Count204       int
	Count301       int
	Count302       int
	Count304       int
	Count400       int
	Count401       int
	Count403       int
	Count404       int
	Count500       int
	Count502       int
	Count503       int
	Count504       int
	Other          int
	AverageLatency int64
	MaxLatency     int64
	MinLatency     int64
}

// ParseURL prepares the URL for the request. It prefixes the URL with
// "http://" or "https://" as appropriate based on the useHTTP flag.
func ParseURL(url string, useHTTP bool) string {
	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
		return url
	}

	if useHTTP {
		return "http://" + url
	}
	return "https://" + url
}

// MakeRequest performs an HTTP GET request to the specified URL.
// It returns an HttpResponse struct filled with response data.
func MakeRequest(useHTTP bool, userAgent, url, headers string) (*HttpResponse, error) {
	var result *HttpResponse
	if userAgent == "" {
		userAgent = "httping"
	}

	tr := &http.Transport{}
	if !useHTTP {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		result.Error = err.Error()
	}

	req.Header.Add("user-agent", userAgent)

	start := time.Now()
	response, err := client.Do(req)
	end := time.Since(start).Milliseconds()
	if err != nil {
		result.Error = err.Error()
	}

	defer func() {
		io.Copy(io.Discard, response.Body)
		response.Body.Close()
	}()

	h := make(map[string]string)
	responseHeaders := strings.Split(headers, ",")

	if len(responseHeaders) > 0 {
		for _, header := range responseHeaders {
			h[header] = response.Header.Get(header)
		}
	}

	result = &HttpResponse{
		Status:          response.StatusCode,
		Host:            response.Header.Get("host"),
		ResponseHeaders: h,
		Latency:         end,
	}

	return result, nil
}

// ParseHeader converts a map of headers into a formatted string.
// This is typically used for displaying the headers in a readable format.
func ParseHeader(m *map[string]string) string {
	var result string
	for k, v := range *m {
		// just making the output pretty
		if len(*m) > 1 {
			result += fmt.Sprintf(" {%s:%s} ", k, v)
		} else {
			result += fmt.Sprintf(" %s:%s ", k, v)
		}
	}
	return result
}

// CalculateStatistics takes a slice of HttpResponse objects and calculates
// aggregated statistics about them, returned as an HTTPStatistics struct.
func CalculateStatistics(responses []*HttpResponse) *HTTPStatistics {
	var stats HTTPStatistics
	var totalLatency int64

	statusCodeCounts := make(map[int]int)

	for _, response := range responses {
		totalLatency += response.Latency
		statusCodeCounts[response.Status]++

		if response.Latency > stats.MaxLatency {
			stats.MaxLatency = response.Latency
		}
		if stats.MinLatency == 0 || response.Latency < stats.MinLatency {
			stats.MinLatency = response.Latency
		}
	}

	stats.Count200 = statusCodeCounts[200]
	stats.Count201 = statusCodeCounts[201]
	stats.Count204 = statusCodeCounts[204]
	stats.Count301 = statusCodeCounts[301]
	stats.Count302 = statusCodeCounts[302]
	stats.Count304 = statusCodeCounts[304]
	stats.Count400 = statusCodeCounts[400]
	stats.Count401 = statusCodeCounts[401]
	stats.Count403 = statusCodeCounts[403]
	stats.Count404 = statusCodeCounts[404]
	stats.Count500 = statusCodeCounts[500]
	stats.Count502 = statusCodeCounts[502]
	stats.Count503 = statusCodeCounts[503]
	stats.Count504 = statusCodeCounts[504]
	stats.Other = 0

	for code, count := range statusCodeCounts {
		switch code {
		case 200, 201, 204, 301, 302, 304, 400, 401, 403, 404, 500, 502, 503, 504:
		default:
			stats.Other += count
		}
	}

	stats.AverageLatency = totalLatency / int64(len(responses))
	return &stats
}

// String provides a formatted string representation of HTTPStatistics,
// making it easy to print the statistics to the console or logs.
func (stats *HTTPStatistics) String() string {
	var builder strings.Builder

	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("AverageLatency: %dms\n", stats.AverageLatency))
	builder.WriteString(fmt.Sprintf("MaxLatency: %dms\n", stats.MaxLatency))
	builder.WriteString(fmt.Sprintf("MinLatency: %dms\n", stats.MinLatency))
	builder.WriteString("\n")

	statusCodes := []struct {
		label string
		count int
	}{
		{"200s", stats.Count200},
		{"201s", stats.Count201},
		{"204s", stats.Count204},
		{"301s", stats.Count301},
		{"302s", stats.Count302},
		{"304s", stats.Count304},
		{"400s", stats.Count400},
		{"401s", stats.Count401},
		{"403s", stats.Count403},
		{"404s", stats.Count404},
		{"500s", stats.Count500},
		{"502s", stats.Count502},
		{"503s", stats.Count503},
		{"504s", stats.Count504},
		{"others", stats.Other},
	}

	for _, sc := range statusCodes {
		builder.WriteString(fmt.Sprintf("Count of %s: %d\n", sc.label, sc.count))
	}

	return builder.String()
}
