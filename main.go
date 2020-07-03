package main

import (
    "fmt"
    "net/http"
    "os"
    "strconv"
    "time"
)

func req(url string, ch chan bool, client *http.Client, keepalive bool) {
    
    //Make one request to url
    req, err := http.NewRequest("GET", url, nil)
    if keepalive {
        req.Header.Set("Connection", "keep-alive")
    }
    resp, err := client.Do(req)
    e := false
    defer resp.Body.Close()
    if err != nil {
        e = true
    }
    ch <- e
    
}

func min(x int, y int) int {
    if x < y {
        return x
    }
    return y
}

func reqs(url string, nreq int, concurrency int, keepalive bool) {
    
    //Make nreq requests with concurrency to url
    tr := &http.Transport{
        DisableKeepAlives: !keepalive,
    }
    clients := make([]*http.Client, concurrency)
    for i := 0; i < concurrency; i++ {
        clients[i] = &http.Client{Transport: tr}
    }
    err := 0
    ch := make(chan bool)
    start := time.Now()
    for i := 0; i < nreq; i += concurrency {
        conc := min(concurrency, nreq-i)
        // Make conc concurrent requests
        for j := 0; j < conc; j++ {
            go req(url, ch, clients[j], keepalive)
        }
        // Wait all requests and collect results
        for j := 0; j < conc; j++ {
            if <- ch {
                err++
            }
        }
    }
    elapsed := float64(time.Since(start).Milliseconds())
    fmt.Println("Time taken for tests: ", elapsed/1000, " seconds")
    fmt.Println("Complete requests: ", nreq-err)
    fmt.Println("Failed requests: ", err)
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
            i++
        } else if arg == "-n" {
            nreq,_ = strconv.Atoi(os.Args[i+1])
            i++
        }
        i++
    }
    
    reqs(url, nreq, concurrency, keepalive)
}
