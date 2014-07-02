package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
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

const (
	STATUS_SUCCESS = iota
	STATUS_FAILED
)

type reqstat struct {
	status int
	time   time.Duration
}

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
	status := reqstat{status: STATUS_FAILED}
	startTime := time.Now()
	resp, err := cli.Do(httpReq)
	if err == nil {
		ioutil.ReadAll(resp.Body) // read entire body
		resp.Body.Close()
		status.time = time.Since(startTime)
		if resp.StatusCode == 200 {
			status.status = STATUS_SUCCESS
		}
	} else {
		status.time = time.Since(startTime)
	}

	statusChan <- status
}

func sendHttpRequests(req *testReq, sendRate, duration, timeout int) *senderstat {

	toDuration := time.Duration(timeout) * time.Second
	client := &http.Client{Timeout: toDuration}
	ticker := time.NewTicker(time.Second / time.Duration(sendRate))
	httpStatusChan := make(chan reqstat)
	doneSendChan := time.After(time.Duration(duration) * time.Second)
	var sstat senderstat
	waitGroup := &sync.WaitGroup{}
	go func() {
		for reqstat := range httpStatusChan {
			if reqstat.status == STATUS_SUCCESS {
				sstat.nsucceded += 1
			} else {
				sstat.nfailed += 1
			}

			sstat.time += reqstat.time
			waitGroup.Done()
		}
	}()

	for {
		select {
		case <-ticker.C:
			waitGroup.Add(1)
			go sendHttpRequest(client, req, httpStatusChan)
		case <-doneSendChan:
			ticker.Stop()
			waitGroup.Wait()
			close(httpStatusChan)
			return &sstat
		}
	}
}

func validateSettings(req *testReq, sendRate, duration, timeout int) error {
	if len(req.url) == 0 {
		return errors.New("URL not specified")
	}

	return nil
}

func main() {
	var req testReq
	sendRate := flag.Int("rate", 50, "send rate (req/s)")
	duration := flag.Int("duration", 3, "send duration (s)")
	timeout := flag.Int("timeout", 20, "request timeout (s)")
	flag.StringVar(&req.method, "method", "GET", "HTTP method")
	flag.StringVar(&req.url, "url", "", "URL")
	flag.StringVar(&req.body, "body", "", "request body")
	flag.StringVar(&req.bodyType, "btype", "application/x-www-form-urlencoded",
		"body type (if body is set)")
	flag.Parse()

	req.method = strings.ToUpper(req.method)
	err := validateSettings(&req, *sendRate, *duration, *timeout)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		res := sendHttpRequests(&req, *sendRate, *duration, *timeout)
		tot := res.nsucceded + res.nfailed
		avg := res.time / time.Duration(tot)
		fmt.Printf("total %v (%v failed) time: %v avg: %v\n", tot, res.nfailed,
			res.time, avg)
	}
}
