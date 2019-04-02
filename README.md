# Echo Server
[![Build Status](https://travis-ci.org/wez470/Echo.svg?branch=master)](https://travis-ci.org/wez470/Echo) 

A simple echo server written in Go

## Building and Running
### Docker
If you are using docker, build the server image and run it with
```
docker build -t echo-server .
docker run -d -p 8080:8080 echo-server
```
Note: `echo-server` can be replaced with any tag you like

###  Without Docker
Run
```
dep ensure
go run .
```
Alternatively, a binary can be built and ran with
```
go build -o echo-server
./echo-server
```
### Auth
To authenticate requests, clients will need to send the `Authorization` header set to `Basic dXNlcjpwYXNzd29yZA==`