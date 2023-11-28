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

## Contributing
Contributions are welcome! Please read our [contribution](CONTRIBUTING.md) guidelines for more information on how to report bugs, suggest enhancements, or submit pull requests.

## To Do
- [x] Add delay option between pings

## Version History
* 1.0.0 - Initial Release

## License
This project is licensed under the MIT License. See the LICENSE.md file for details.
