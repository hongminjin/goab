package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "time"
)

// Make one request to url using client and return if there is an error
func req(url string, client *http.Client, keepalive bool) bool {

    req, err := http.NewRequest("GET", url, nil)
    if keepalive {
        req.Header.Set("Connection", "keep-alive")
    }
    resp, err := client.Do(req)
    e := false
    defer resp.Body.Close()
    ioutil.ReadAll(resp.Body)
    if err != nil {
        e = true
    }
    return e

}

// Make requests to url and send the number of errors to channel cerr
func work(url string, keepalive bool, cerr chan int, ch chan bool) {

    tr := &http.Transport{
        DisableKeepAlives: !keepalive,
    }
    client := &http.Client{Transport: tr} // client to use
    err := 0 // total number of errors

    // receive jobs from ch
    v, ok := <- ch
    for ; v && ok; {
        e := req(url, client, keepalive)
        if e {
            err++
        }
        v, ok = <- ch
    }
    cerr <- err

}

func min(x int, y int) int {
    if x < y {
        return x
    }
    return y
}

// Make nreq requests with concurrency to url and print results
func reqs(url string, nreq int, concurrency int, keepalive bool) {

    err := 0 // total number of errors
    cerr := make(chan int) // channel to collect errors
    ch := make(chan bool, concurrency) // channel to assign jobs
    conc := min(nreq, concurrency) // number of goroutines

    // create conc goroutines to do the job
    for i := 0; i < conc; i++ {
        go work(url, keepalive, cerr, ch)
    }

    start := time.Now()
    // create jobs
    for i := 0; i < nreq; i++ {
        ch <- true
    }
    for i := 0; i < conc; i++ {
        ch <- false
    }
    close(ch)
    // wait all goroutines and collect results by channel cerr
    for i := 0; i < conc; i++ {
        e := <- cerr
        err += e
    }
    elapsed := float64(time.Since(start).Milliseconds())

    fmt.Println("Time taken for tests: ", elapsed/1000, " seconds")
    fmt.Println("Complete requests: ", nreq-err)
    fmt.Println("Failed requests: ", err)
    fmt.Println("TPS: ", float64(nreq)/(elapsed/1000), " [#/sec] (mean)")
    fmt.Println("Time per request: ", elapsed*float64(conc)/float64(nreq), " [ms] (mean)")
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
