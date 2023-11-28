package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/rnemeth90/httping"
	"github.com/spf13/pflag"
)

type config struct {
	url     string
	useHTTP bool
	count   int
	headers string
}

var (
	url     string
	useHTTP bool
	count   int
	headers string
	sleep   int64
)

func init() {
	pflag.StringVar(&url, "url", "", "Specify the URL to ping. (required)")
	pflag.BoolVar(&useHTTP, "insecure", false, "Use HTTP instead of HTTPS. By default, HTTPS is used.")
	pflag.StringVar(&headers, "headers", "", "A comma-separated list of response headers to include in the output.")
	pflag.IntVar(&count, "count", 4, "Set the number of pings to send. Default is 4.")
	pflag.Int64Var(&sleep, "sleep", 0, "Set the delay (in seconds) between successive pings. Default is 0 (no delay).")
	pflag.ErrHelp = nil
	pflag.Usage = usage
}

func usage() {
	fmt.Println("httping: A tool to 'ping' a web server and display response statistics.")
	fmt.Println("\nUsage:")
	fmt.Println("  httping [OPTIONS] --url URL")
	fmt.Println("\nExamples:")
	fmt.Println("  httping --url www.google.com")
	fmt.Println("  httping --url www.google.com --insecure --count 10")
	fmt.Println("  httping --url www.google.com --count 100 --headers Content-Type,Server")
	fmt.Println("  httping --url www.google.com --sleep 10")

	fmt.Println("Options:")
	pflag.PrintDefaults()
}

func main() {
	pflag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(pflag.Args()) == 0 {
		usage()
		os.Exit(0)
	}

	osArgFlags := pflag.Args()
	if len(osArgFlags) > 0 {
		url = httping.ParseURL(os.Args[1], useHTTP)
	} else {
		url = httping.ParseURL(url, useHTTP)
	}

	config := config{
		url:     url,
		useHTTP: useHTTP,
		count:   count,
		headers: headers,
	}

	go func() {
		osChan := make(chan os.Signal, 1)
		signal.Notify(osChan, os.Interrupt, syscall.SIGTERM)

		<-osChan
		cancel()
	}()

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)

	if err := run(ctx, config, tw); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, config config, writer io.Writer) error {
	var count int
	var respForStats []*httping.HttpResponse

	// check if the writer is a tabwriter
	tw, ok := writer.(*tabwriter.Writer)

	fmt.Fprintln(writer, "Time\tCount\tUrl\tResult\tTime\tHeaders")
	fmt.Fprintln(writer, "-----\t-----\t---\t------\t----\t-------")

	for i := 1; i <= config.count; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Operation cancelled by user")
			stats := httping.CalculateStatistics(respForStats)
			fmt.Printf("Total Requests: %d\n", count)
			fmt.Println(stats.String())
			return ctx.Err()
		default:
			response, err := httping.MakeRequest(config.useHTTP, config.url, config.headers)
			if err != nil {
				return err
			}
			respForStats = append(respForStats, response)

			headerValues := httping.ParseHeader(&response.ResponseHeaders)

			hs := *headerValues
			fmt.Fprintf(writer, "[ %v ]\t[ %d ]\t[ %s ]\t[ %d ]\t[ %dms ]\t[ %s ]\n", time.Now().Format(time.RFC3339), i, config.url, response.Status, response.Latency, hs)
			count++

			if ok {
				tw.Flush()
			}
			time.Sleep(time.Second * time.Duration(sleep))
		}
	}

	stats := httping.CalculateStatistics(respForStats)
	fmt.Println()
	fmt.Printf("Total Requests: %d\n", count)
	fmt.Println(stats.String())
	return nil
}
