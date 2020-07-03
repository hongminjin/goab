package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    
    keepalive := false // HTTP keepalive
    nreq := 1 // Number of requests
    concurrency := 1 // Number of multiple requests to make at a time
    
    for i := 1; i < len(os.Args); {
        arg := os.Args[i]
        if arg == "-k" {
            keepalive = true
        } else if arg == "-c" {
            concurrency,_ = strconv.Atoi(os.Args[i+1])
            i = i+1
        } else if arg == "-n" {
            nreq,_ = strconv.Atoi(os.Args[i+1])
            i = i+1
        }
        i = i+1
    }
    fmt.Println(keepalive, nreq, concurrency)
}
