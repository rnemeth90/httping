package httping

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type HttpResponse struct {
	Status          int
	Host            string
	ResponseHeaders map[string]string
	Latency         int64
}

// HTTPStatistics defines a list of integers for keeping track of the
// total number of HTTP response codes for a given response code. We
// only account for the most common response codes.
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

// ParseURL will parse a URL from a string
func ParseURL(url string, useHTTP bool) string {
	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
		return url
	}

	if useHTTP {
		return "http://" + url
	}
	return "https://" + url
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

	result = &HttpResponse{
		Status:          response.StatusCode,
		Host:            response.Header.Get("host"),
		ResponseHeaders: h,
		Latency:         end,
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

func CalculateStatistics(responses []*HttpResponse) *HTTPStatistics {
	var stats HTTPStatistics
	var totalLatency int64

	for _, l := range responses {
		totalLatency += l.Latency

		switch l.Status {
		case 200:
			stats.Count200++
		case 201:
			stats.Count201++
		case 204:
			stats.Count204++
		case 301:
			stats.Count301++
		case 302:
			stats.Count302++
		case 304:
			stats.Count304++
		case 400:
			stats.Count400++
		case 401:
			stats.Count401++
		case 403:
			stats.Count403++
		case 404:
			stats.Count404++
		case 500:
			stats.Count500++
		case 502:
			stats.Count502++
		case 503:
			stats.Count503++
		case 504:
			stats.Count504++
		default:
			stats.Other++
		}
	}

	var max int64
	for i, v := range responses {
		if i == 0 || v.Latency > max {
			max = v.Latency
		}
	}

	var min int64
	for i, v := range responses {
		if i == 0 || v.Latency < min {
			min = v.Latency
		}
	}

	stats.AverageLatency = totalLatency / int64(len(responses))
	stats.MaxLatency = max
	stats.MinLatency = min
	return &stats
}

func (stats *HTTPStatistics) String() string {
	var b bytes.Buffer

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("AverageLatency: %dms\n", stats.AverageLatency))
	b.WriteString(fmt.Sprintf("MaxLatency: %dms\n", stats.MaxLatency))
	b.WriteString(fmt.Sprintf("MinLatency: %dms\n", stats.MinLatency))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Count of 200s: %d\n", stats.Count200))
	b.WriteString(fmt.Sprintf("Count of 201s: %d\n", stats.Count201))
	b.WriteString(fmt.Sprintf("Count of 204s: %d\n", stats.Count204))
	b.WriteString(fmt.Sprintf("Count of 301s: %d\n", stats.Count301))
	b.WriteString(fmt.Sprintf("Count of 302s: %d\n", stats.Count302))
	b.WriteString(fmt.Sprintf("Count of 304s: %d\n", stats.Count304))
	b.WriteString(fmt.Sprintf("Count of 400s: %d\n", stats.Count400))
	b.WriteString(fmt.Sprintf("Count of 401s: %d\n", stats.Count401))
	b.WriteString(fmt.Sprintf("Count of 403s: %d\n", stats.Count403))
	b.WriteString(fmt.Sprintf("Count of 404s: %d\n", stats.Count404))
	b.WriteString(fmt.Sprintf("Count of 500s: %d\n", stats.Count500))
	b.WriteString(fmt.Sprintf("Count of 502s: %d\n", stats.Count502))
	b.WriteString(fmt.Sprintf("Count of 503s: %d\n", stats.Count503))
	b.WriteString(fmt.Sprintf("Count of 504s: %d\n", stats.Count504))
	b.WriteString(fmt.Sprintf("Count of others: %d\n", stats.Other))

	return b.String()
}
