package main

import (
    "fmt"
    "net/http"
    "os"
    "strconv"
    "time"
)

func req(url string, client *http.Client, keepalive bool) bool {
    
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
    return e
    
}

func work(url string, keepalive bool, ch chan int, n int) {
    tr := &http.Transport{
        DisableKeepAlives: !keepalive,
    }
    client := &http.Client{Transport: tr}
    err := 0
    for i := 0; i < n; i++ {
        e := req(url, client, keepalive)
        if e {
            err++
        }
    }
    ch <- err
}

/*func min(x int, y int) int {
    if x < y {
        return x
    }
    return y
}*/

func reqs(url string, nreq int, concurrency int, keepalive bool) {
    
    //Make nreq requests with concurrency to url
    err := 0
    ch := make(chan int)
    n := make([]int, concurrency)
    for i := 0; i < concurrency; i++ {
        j := 0
        if i < nreq%concurrency {
            j = 1
        }
        n[i] = nreq/concurrency + j
    }
    start := time.Now()
    for i := 0; i < concurrency; i++ {
        go work(url, keepalive, ch, n[i])
    }
    for i := 0; i < concurrency; i++ {
        e := <- ch
        err += e
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
