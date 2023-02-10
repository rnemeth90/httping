# httping [![build-release-binary](https://github.com/rnemeth90/httping/actions/workflows/build.yaml/badge.svg)](https://github.com/rnemeth90/httping/actions/workflows/build.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/rnemeth90/httping/)](https://goreportcard.com/report/github.com/rnemeth90/httping/)
## Description
httping is a simple tool for 'pinging' a web server. You can specify the count of pings, as well as some HTTP headers to return the value of. After all pings have completed, some statistics are printed to stdout.

## Getting Started
To get started, simply download the latest release using the link below and execute the program with `--help`
```
gopher$ httping --help
httping

Usage:
  httping --url www.google.com
  httping --url www.google.com --c 100 --headers server

Options:
      --c int            number of pings to send (default 4)
      --headers string   comma delimited list of response headers to output
      --insecure         use http instead of https
      --url string       the url to ping
```

### Dependencies
* to build yourself, you must have Go v1.13+ installed

### Installing
Simply download the latest release [here](https://github.com/rnemeth90/httping/releases)

### Executing program
```
gopher$ httping --url www.google.com --c 100 --headers server
Time                            Count   Url                             Result  Time            Headers
-----                           -----   ---                             ------  ----            -------
[ 2023-02-10T13:14:45-05:00 ]   [ 1 ]   [ https://www.google.com ]      [ 200 ] [ 228ms ]       [  server:gws  ]
[ 2023-02-10T13:14:46-05:00 ]   [ 2 ]   [ https://www.google.com ]      [ 200 ] [ 199ms ]       [  server:gws  ]
[ 2023-02-10T13:14:46-05:00 ]   [ 3 ]   [ https://www.google.com ]      [ 200 ] [ 234ms ]       [  server:gws  ]
[ 2023-02-10T13:14:46-05:00 ]   [ 4 ]   [ https://www.google.com ]      [ 200 ] [ 171ms ]       [  server:gws  ]
[ 2023-02-10T13:14:46-05:00 ]   [ 5 ]   [ https://www.google.com ]      [ 200 ] [ 184ms ]       [  server:gws  ]
[ 2023-02-10T13:14:46-05:00 ]   [ 6 ]   [ https://www.google.com ]      [ 200 ] [ 171ms ]       [  server:gws  ]
[ 2023-02-10T13:14:47-05:00 ]   [ 7 ]   [ https://www.google.com ]      [ 200 ] [ 176ms ]       [  server:gws  ]
[ 2023-02-10T13:14:47-05:00 ]   [ 8 ]   [ https://www.google.com ]      [ 200 ] [ 172ms ]       [  server:gws  ]
^C
Total Requests: 8

AverageLatency: 191ms
MaxLatency: 234ms
MinLatency: 171ms

Count of 200s: 8
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

## Help
If you need help or have suggestions/feature requests, please submit an issue.

## To Do
- [x] replace flag with pflag
- [x] add statistics
- [x] add stringer for httpstatistics
- [x] add tests
- [x] finish readme
- [x] allow host to be input as arg 1

## Version History
* 1.0.0
    * Initial Release

## License
This project is licensed under the MIT License - see the LICENSE.md file for details
