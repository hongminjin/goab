# goab
Simple implementation of Apache Benchmark with Go

## Description

* documentation.pdf: documentation of the project
* main.go: main program
* server.go: to run a HTTP server in the localhost

## Instructions

To create executable of the main program and server, use:
```
go build -o goab main.go
go build -o server server.go
```

## Usage

```
./goab [parameters] url
```
Where `url` is the url to test, available parameters are:
* `-k` to enable HTTP KeepAlive
* `-n n` to make n requests
* `-c c` to specify concurrency (number of multiple requests to perform at a time)

To test localhost, run
```
./server
```
and
```
./goab [parameters] http://localhost:8080/
```
