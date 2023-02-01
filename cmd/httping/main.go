package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/rnemeth90/httping"
)

type config struct {
	url string
	useHTTP bool
	count int
	headers string
}

var(
	url string
	useHTTP bool
	co int
	headers string
)

func init() {
	flag.StringVar(&url, "url", "", "the url to ping")
	flag.BoolVar(&useHTTP, "insecure", false, "use http instead of https")
	flag.StringVar(&headers, "headers", "", "comma delimited list of response headers to output")
	flag.IntVar(&co, "c", 4, "number of pings to send")
	flag.Usage = usage
}

func usage() {
	fmt.Println(os.Args[0])

	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  httping -u www.google.com\n\n")

	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	url = httping.ParseURL(url, useHTTP)

	c := config {
		url: url,
		useHTTP: useHTTP,
		count: co,
		headers: headers,
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)

	if err := run(c, tw); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(c config, writer io.Writer) error {
	var count int

	tw, ok := writer.(*tabwriter.Writer)

	// tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "Time\tCount\tUrl\tResult\tTime\tHeaders")
	fmt.Fprintln(writer, "-----\t-----\t---\t------\t----\t-------")

	osChan := make(chan os.Signal)
	signal.Notify(osChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-osChan
		fmt.Println()
		fmt.Printf("Total Requests: %d\n", count)
		os.Exit(0)
	}()


	for i := 1; i <= c.count; i++ {
		response, err := httping.MakeRequest(c.url, c.headers)
		if err != nil {
			return err
		}

		s := response.Status
		l := response.Latency

		headerValues := httping.ParseMap(&response.ResponseHeaders)

		hs := *headerValues
		fmt.Fprintf(writer, "[ %v ]\t[ %d ]\t[ %s ]\t[ %s ]\t[ %dms ]\t[ %s ]\n", time.Now().Format(time.RFC3339), i, c.url, s, l, hs)
		count++

		if ok {
			tw.Flush()
		}
	}

	return nil
}