package main

import (
    "fmt"
    "net/http"
    "os"
    "sort"
    "strconv"
    "time"
)

func req(url string, ch chan int) {
    
    //Make one request to url and counting milliseconds
    start := time.Now()
    resp, err := http.Get(url)
    t := int(time.Since(start).Milliseconds())
    defer resp.Body.Close()
    if err != nil {
        t = -1
    }
    ch <- t
    
}

func min(x int, y int) int {
    if x < y {
        return x
    }
    return y
}

func reqs(url string, nreq int, concurrency int) {
    
    //Make nreq requests with concurrency to url
    timereq := make([]int, nreq)
    ch := make(chan int)
    start := time.Now()
    for i := 0; i < nreq; {
        conc := min(concurrency, nreq-i)
        // Make conc concurrent requests
        for j := 0; j < conc; j++ {
            go req(url, ch)
        }
        // Wait all requests and collect results
        for j := 0; j < conc; j++ {
            timereq[i] = <- ch
            i++
        }
    }
    elapsed := float64(time.Since(start).Milliseconds())
    sort.Ints(timereq[:])
    fmt.Println("Time taken for tests: ", elapsed/1000, " seconds")
    fmt.Println("TPS: ", float64(nreq)/(elapsed/1000), " [#/sec] (mean)")
    fmt.Println("Time per request: ", elapsed*float64(concurrency)/float64(nreq), " [ms] (mean)")
    fmt.Println("Time per request: ", elapsed/float64(nreq), "[ms] (mean, across all concurrent requests)")
}

func main() {
    
    keepalive := false // HTTP keepalive
    nreq := 1 // Number of requests
    concurrency := 1 // Number of multiple requests to make at a time
    url := os.Args[len(os.Args)-1] // Request url
    
    for i := 1; i < len(os.Args)-1; {
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
    
    if keepalive {
        fmt.Println("keepalive")
    }
    
    reqs(url, nreq, concurrency)
}
