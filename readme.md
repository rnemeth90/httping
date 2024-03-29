# httping [![build-release-binary](https://github.com/rnemeth90/httping/actions/workflows/build.yaml/badge.svg)](https://github.com/rnemeth90/httping/actions/workflows/build.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/rnemeth90/httping)](https://goreportcard.com/report/github.com/rnemeth90/httping)

## Overview
`httping` is a command-line utility designed to perform HTTP pings to web servers, providing key metrics like response times and status code counts. It's a handy tool for monitoring and testing website performance.

## Features
- Ping a web server with a specified URL.
- Configure the number of pings.
- Choose to include specific HTTP headers in the output.
- View summarized statistics after pinging, including latency metrics and response code distribution.

### Dependencies
* to build yourself, you must have Go v1.13+ installed

### Installing
Download the latest release from [here](https://github.com/rnemeth90/httping/releases).

### Usage
To see all available options:
```
$ httping --help

httping: A tool to 'ping' a web server and display response statistics.

Usage:
  httping [OPTIONS] --url URL

Examples:
  httping --url www.google.com
  httping --url www.google.com --insecure --count 10
  httping --url www.google.com --count 100 --headers Content-Type,Server
  httping --url www.google.com --sleep 10
Options:
      --count int        Set the number of pings to send. Default is 4. (default 4)
      --headers string   A comma-separated list of response headers to include in the output.
      --insecure         Use HTTP instead of HTTPS. By default, HTTPS is used.
      --sleep int        Set the delay (in seconds) between successive pings. Default is 0 (no delay).
      --url string       Specify the URL to ping. (required)
```

Basic usage:
```
$ httping --url www.google.com
```

With custom options:
```
$ httping --url www.google.com --c 100 --headers server
```

### Example Output
Basic usage:
```
$ httping www.google.com
Time                            Count   Url                             Result  Time            Headers
-----                           -----   ---                             ------  ----            -------
[ 2023-11-28T11:44:05-05:00 ]   [ 1 ]   [ https://www.google.com ]      [ 200 ] [ 312ms ]       [  :  ]
[ 2023-11-28T11:44:05-05:00 ]   [ 2 ]   [ https://www.google.com ]      [ 200 ] [ 186ms ]       [  :  ]
[ 2023-11-28T11:44:05-05:00 ]   [ 3 ]   [ https://www.google.com ]      [ 200 ] [ 171ms ]       [  :  ]
[ 2023-11-28T11:44:05-05:00 ]   [ 4 ]   [ https://www.google.com ]      [ 200 ] [ 206ms ]       [  :  ]

Total Requests: 4

AverageLatency: 218ms
MaxLatency: 312ms
MinLatency: 171ms

Count of 200s: 4
Count of 201s: 0
Count of 204s: 0
Count of 301s: 0
Count of 302s: 0
Count of 304s: 0
Count of 400s: 0
Count of 401s: 0
Count of 403s: 0
Count of 404s: 0
Count of 500s: 0
Count of 502s: 0
Count of 503s: 0
Count of 504s: 0
Count of others: 0
```

With custom options:
```
 $ httping --url www.google.com --count 10 --headers Content-Type,Server
Time                            Count   Url                             Result  Time            Headers
-----                           -----   ---                             ------  ----            -------
[ 2023-11-28T11:48:25-05:00 ]   [ 1 ]   [ https://www.google.com ]      [ 200 ] [ 211ms ]       [  {Content-Type:text/html; charset=ISO-8859-1}  {Server:gws}  ]
[ 2023-11-28T11:48:25-05:00 ]   [ 2 ]   [ https://www.google.com ]      [ 200 ] [ 182ms ]       [  {Server:gws}  {Content-Type:text/html; charset=ISO-8859-1}  ]
[ 2023-11-28T11:48:25-05:00 ]   [ 3 ]   [ https://www.google.com ]      [ 200 ] [ 183ms ]       [  {Content-Type:text/html; charset=ISO-8859-1}  {Server:gws}  ]
[ 2023-11-28T11:48:25-05:00 ]   [ 4 ]   [ https://www.google.com ]      [ 200 ] [ 176ms ]       [  {Content-Type:text/html; charset=ISO-8859-1}  {Server:gws}  ]
[ 2023-11-28T11:48:25-05:00 ]   [ 5 ]   [ https://www.google.com ]      [ 200 ] [ 189ms ]       [  {Content-Type:text/html; charset=ISO-8859-1}  {Server:gws}  ]
[ 2023-11-28T11:48:26-05:00 ]   [ 6 ]   [ https://www.google.com ]      [ 200 ] [ 206ms ]       [  {Content-Type:text/html; charset=ISO-8859-1}  {Server:gws}  ]
[ 2023-11-28T11:48:26-05:00 ]   [ 7 ]   [ https://www.google.com ]      [ 200 ] [ 170ms ]       [  {Content-Type:text/html; charset=ISO-8859-1}  {Server:gws}  ]
[ 2023-11-28T11:48:26-05:00 ]   [ 8 ]   [ https://www.google.com ]      [ 200 ] [ 171ms ]       [  {Content-Type:text/html; charset=ISO-8859-1}  {Server:gws}  ]
[ 2023-11-28T11:48:26-05:00 ]   [ 9 ]   [ https://www.google.com ]      [ 200 ] [ 177ms ]       [  {Content-Type:text/html; charset=ISO-8859-1}  {Server:gws}  ]
[ 2023-11-28T11:48:26-05:00 ]   [ 10 ]  [ https://www.google.com ]      [ 200 ] [ 209ms ]       [  {Content-Type:text/html; charset=ISO-8859-1}  {Server:gws}  ]

Total Requests: 10

AverageLatency: 187ms
MaxLatency: 211ms
MinLatency: 170ms

Count of 200s: 10
Count of 201s: 0
Count of 204s: 0
Count of 301s: 0
Count of 302s: 0
Count of 304s: 0
Count of 400s: 0
Count of 401s: 0
Count of 403s: 0
Count of 404s: 0
Count of 500s: 0
Count of 502s: 0
Count of 503s: 0
Count of 504s: 0
Count of others: 0
```

## Contributing
Contributions are welcome! Please read our [contribution](CONTRIBUTING.md) guidelines for more information on how to report bugs, suggest enhancements, or submit pull requests.

## To Do
- [x] Add delay option between pings
- [x] Improve help messages and usage()
- [x] Ensure POSIX compliance
- [x] The current approach of using a goroutine for signal handling can be improved for clarity and control. Consider moving signal handling to the main function or using a context with cancel functionality.
- [x] remove 'context canceled' and 'operation cancelled by user' when using SIGTERM
- [x] When an error occurs in MakeRequest, it immediately returns from the function. Consider adding logic to handle partial results and provide a summary of successes and failures.
- [x] The call to httping.ParseHeader and dereferencing the result (*headerValues) can be optimized. Maybe modify ParseHeader to directly return a string.
- [x] The check if ok { tw.Flush() } is done in every iteration. If you are sure that the writer will always be a tabwriter.Writer, this check can be done once before the loop.
- [x] The current implementation always sleeps after each request, even if sleep is 0. You can optimize this by adding a conditional check to avoid unnecessary sleeping.
- [x] The final statistics and request count are printed at the end of the function. It's good practice to also handle situations where the loop might exit unexpectedly.
- [x] Add documentation comments for all exported functions and entities
- [ ] The time format in the output is hard-coded. Consider making this format configurable via command-line arguments or configuration files.
- [x] Add tests. Moving the logic that processes individual responses into a separate function can make it easier to write unit tests.

## Version History
* 1.0.0 - Initial Release
* 1.1.2 - Fixes, enhancements, refactoring, tests, etc.

## License
This project is licensed under the MIT License. See the LICENSE.md file for details.
