## Problem

black box testing against live systems can cause high load on the
tested system, resulting in loss of availability. 
 
It is often assumed that the raw number of request sent to a system
correlates to the load (CPU, I/O wait) on this system over time. While this 
may or may not be true, limiting the number of requests sent per unit of time 
to a system can result in limiting the overall load on that system. Limiting 
the number of requests sent per unit of time will also increase scan time, 
which can result in decreased performance at the receiving end of the scan 
over a longer period of time than a faster scan would, and can also affect the 
test results negatively (missed tests, &c).

One of the ways to try and achieve a proper send rate is to assume that 
there's a dependence between the number of requests sent per unit of time, and
response times for those requests.  The most common assumption is that an 
increase in the number of sent requests will lead to an increase in response 
time. Based on this assumption, some tools may try to adapt the send rate 
based on response time measurements. 

We will explore that assumption in the environment of web application
scanning tools used over the Internet.

## http-stress-test.go test scenario

```
$ ./http-stress-test -h
Usage of http-stress-test:
  -body="": request body
  -btype="application/x-www-form-urlencoded": body type (if body is set)
  -duration=3: send duration (s)
  -method="GET": HTTP method
  -rate=50: send rate (req/s)
  -timeout=20: request timeout (s)
  -url="": URL
```

We're using Go's net/http DefaultClient as a reference client with
no extra options. The complete response body is read. DefaultClient supports 
compression and keep-alive functionality. This is reasonable functionality to 
have in an HTTP based scanner.

Difference in average response times between requests to different resources 
is expected. foo.php might execute more code and perform more DB interaction 
server side that bar.php. Therefor, we limit our tests to one resource. 

The tests are performed over the Internet, no WiFi, against a server with
no other significant traffic than the test traffic itself. 


## http-stress-test.go measurements


