package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type testReq struct {
	method   string
	url      string
	bodyType string
	body     string
}

type reqstat int

const (
	STATUS_SUCCESS = iota
	STATUS_FAILED
)

type senderstat struct {
	nsucceded int
	nfailed   int
	time      time.Duration
}

func sendHttpRequest(cli *http.Client, req *testReq, statusChan chan reqstat) {
	var httpReq *http.Request
	if len(req.body) > 0 {
		body := bytes.NewReader([]byte(req.body))
		httpReq, _ = http.NewRequest(req.method, req.url, body)
		httpReq.Header.Set("Content-Type", req.bodyType)
	} else {
		httpReq, _ = http.NewRequest(req.method, req.url, nil)
	}

	httpReq.Header.Set("User-Agent", "HTTPStressTester/1.0")
	resp, err := cli.Do(httpReq)
	if err == nil && resp.StatusCode == 200 {
		statusChan <- STATUS_SUCCESS
	} else {
		statusChan <- STATUS_FAILED
	}
}

func startHttpSender(req *testReq, sendRate, duration int) *senderstat {

	client := &http.Client{}
	ticker := time.NewTicker(time.Second / time.Duration(sendRate))
	httpStatusChan := make(chan reqstat)
	doneSendChan := time.After(time.Duration(duration) * time.Second)
	var sstat senderstat
	waitGroup := &sync.WaitGroup{}
	startTime := time.Now()
	go func() {
		for rstat := range httpStatusChan {
			if rstat == STATUS_SUCCESS {
				sstat.nsucceded += 1
			} else {
				sstat.nfailed += 1
			}
			waitGroup.Done()
		}
	}()

	for {
		select {
		case <-ticker.C:
			go sendHttpRequest(client, req, httpStatusChan)
			waitGroup.Add(1)
		case <-doneSendChan:
			ticker.Stop()
			waitGroup.Wait()
			close(httpStatusChan)
			sstat.time = time.Since(startTime)
			return &sstat
		}
	}
}

func validateSettings(req *testReq, sendRate, duration int) error {
	if len(req.url) == 0 {
		return errors.New("URL not specified")
	}

	return nil
}

func main() {
	var req testReq
	sendRate := flag.Int("rate", 50, "send rate (req/s)")
	duration := flag.Int("duration", 3, "send duration (s)")
	flag.StringVar(&req.method, "method", "GET", "HTTP method")
	flag.StringVar(&req.url, "url", "", "URL")
	flag.StringVar(&req.body, "body", "", "request body")
	flag.StringVar(&req.bodyType, "btype", "application/x-www-form-urlencoded",
		"body type (if body is set)")
	flag.Parse()

	req.method = strings.ToUpper(req.method)

	err := validateSettings(&req, *sendRate, *duration)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		res := startHttpSender(&req, *sendRate, *duration)
		tot := res.nsucceded + res.nfailed
		avg := res.time / time.Duration(tot)
		fmt.Printf("total %v (%v failed) time: %v avg: %v\n", tot, res.nfailed,
			res.time, avg)
	}
}
