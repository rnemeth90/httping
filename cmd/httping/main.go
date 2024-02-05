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

// config holds the configuration parameters for httping execution.
// It includes the target URL, HTTP/HTTPS usage preference, count of pings,
// headers to be retrieved, and sleep duration between pings.
type config struct {
	url       string
	useHTTP   bool
	count     int
	headers   string
	sleep     int64
	useragent string
}

var (
	url       string
	useHTTP   bool
	count     int
	headers   string
	sleep     int64
	useragent string
)

// init is an initialization function that sets up command-line flags
// and arguments for the httping tool. It defines the flags, default values,
// and help messages for each command-line option. The usage function is also
// set here to provide custom help text.
func init() {
	pflag.StringVar(&url, "url", "", "Specify the URL to ping. (required)")
	pflag.BoolVar(&useHTTP, "insecure", false, "Use HTTP instead of HTTPS. By default, HTTPS is used.")
	pflag.StringVar(&headers, "headers", "", "A comma-separated list of response headers to include in the output.")
	pflag.StringVar(&useragent, "user-agent", "", "The user-agent value to include in the request headers.")
	pflag.IntVar(&count, "count", 4, "Set the number of pings to send. Default is 4.")
	pflag.Int64Var(&sleep, "sleep", 0, "Set the delay (in seconds) between successive pings. Default is 0 (no delay).")
	pflag.Usage = usage
}

// usage prints the help text to the console. It provides a detailed
// description of how to use the httping tool, including its syntax,
// available options, and examples of common usages. This function is
// designed to guide the user in effectively utilizing the tool.
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

// main is main :)
func main() {
	pflag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(pflag.Args()) == 0 && url == "" {
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
		url:       url,
		useHTTP:   useHTTP,
		count:     count,
		headers:   headers,
		sleep:     sleep,
		useragent: useragent,
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

// run executes the httping process based on the provided configuration and context.
// It handles pinging the specified URL, collecting response data, and writing the
// output to the provided writer. The function respects context cancellation,
// allowing for graceful termination. It accumulates statistics throughout the
// execution and prints a summary at the end or upon early termination.
func run(ctx context.Context, config config, writer io.Writer) error {
	var count int
	var respForStats []*httping.HttpResponse

	defer func() {
		if len(respForStats) > 0 {
			stats := httping.CalculateStatistics(respForStats)
			fmt.Printf("Total Requests: %d\n", count)
			fmt.Println(stats.String())
		}
	}()

	// check if the writer is a tabwriter
	tw, isTabWriter := writer.(*tabwriter.Writer)
	defer tw.Flush()

	fmt.Fprintln(writer, "Time\tCount\tUrl\tResult\tTime\tHeaders")
	fmt.Fprintln(writer, "-----\t-----\t---\t------\t----\t-------")

	for i := 1; i <= config.count; i++ {
		select {
		case <-ctx.Done():
			stats := httping.CalculateStatistics(respForStats)
			fmt.Printf("Total Requests: %d\n", count)
			fmt.Println(stats.String())
			return ctx.Err()
		default:
			response, err := httping.MakeRequest(config.useHTTP, config.useragent, config.url, config.headers)
			if err != nil {
				return err
			}

			respForStats = append(respForStats, response)

			headerValues := httping.ParseHeader(&response.ResponseHeaders)
			fmt.Fprintf(writer, "[ %v ]\t[ %d ]\t[ %s ]\t[ %d ]\t[ %dms ]\t[ %s ]\n", time.Now().Format(time.RFC3339), i, config.url, response.Status, response.Latency, headerValues)

			count++

			// we flush the tab writer in the loop for scrolling "live" output
			if isTabWriter {
				tw.Flush()
			}

			if config.sleep != 0 {
				time.Sleep(time.Second * time.Duration(config.sleep))
			}
		}
	}

	return nil
}
