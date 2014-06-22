package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
)

const NCONCURRENTMAX = 20
const NREQS = 5
const DSTURL = "http://www.aftonbladet.se/"

type httpStats struct {
	nsucceded int
	nfailed   int
}

func sendHttpRequests(url string, nreqs int, statsChan chan httpStats) {
	var stats httpStats
	for i := 0; i < nreqs; i++ {
		resp, err := http.Get(url)
		if err == nil {
			if resp.StatusCode == 200 {
				stats.nsucceded += 1
			} else {
				stats.nfailed += 1
			}
			
			// read entire response
			ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		} else {
			stats.nfailed += 1
		}
	}

	statsChan <- stats
}

func dispatchHttpRequesters(nreqs, nconcurr int, dsturl string, stats *httpStats) {
	requesterChan := make(chan httpStats)
	for i := 0; i < nconcurr; i++ {
		go sendHttpRequests(dsturl, nreqs, requesterChan)
	}

	for i := 0; i < nconcurr; i++ {
		rstat := <-requesterChan
		if stats != nil {
			stats.nsucceded += rstat.nsucceded
			stats.nfailed += rstat.nfailed
		}
	}
}

func main() {
	var stats httpStats

	// set up connection cache, &c	
	dispatchHttpRequesters(NREQS, 1, DSTURL, nil)

	for i := 1; i < NCONCURRENTMAX; i++ {
		stats.nsucceded = 0
		stats.nfailed = 0
		startTime := time.Now()
		dispatchHttpRequesters(NREQS, i, DSTURL, &stats)
		duration := time.Since(startTime)
		avg := duration.Nanoseconds() / (int64(stats.nsucceded) + int64(stats.nfailed)) / 1000000;
		fmt.Printf("nconcurrent: %v nreqs: %v succeded: %v failed: %v duration: %v avg: %vms\n",
			i, NREQS, stats.nsucceded, stats.nfailed, duration, avg)
		time.Sleep(5 * time.Second)
	}
}
