## Problem

black box testing against live systems can cause high load on the
tested system, resulting in loss of availability. 
 
It is often assumed that the raw number of request sent to a system
correlates to the load (CPU, I/O wait) on this system over time. While this 
may or may not be true, limiting the number of requests sent per unit of time 
to a system can result in limiting the overall load on that system. Limiting 
the number of requests sent per unit of time will also increase scan time, 
which can result in decreased performance over a longer period of time and can 
also affect the test results negatively (expired sessions, &c).

One of the ways to try and achieve a proper send rate is to assume that 
there's a dependence between the number of requests sent per unit of time, and
response times for those requests.  The most common assumption is that an 
increase in the number of sent requests will lead to an increase in response 
time. Based on this assumption, some tools may try to adapt the send rate 
based on response time measurements. 

We will explore that assumption in the environment of web application
scanning tools used over the Internet.

## Sender model for weblat.go

We're using Go's net/http DefaultClient as a reference client with
no extra options. The complete response body is read. DefaultClient supports 
compression and keep-alive functionality. This is reasonable functionality to 
have in an HTTP based scanner.

The sender model will have a number of concurrent processes sending an HTTP 
request and getting a response in sequence, NREQS time in a row. The number of
concurrent processes (implemented as goroutines) will increase from 1 to 
NCONCURRENTMAX during a test run. When all the concurrent processes have sent
NREQS HTTP requests and received their responses or timed out, the test will 
be suspended for five seconds until the next round will begin, this time with
one more concurrent process.  The number of requests sent by each 
concurrent process is given by the constant NREQS. All HTTP requests are 
directed against one resource, DSTURL. 

Difference in average response times between requests to different resources 
is expected. foo.php might execute more code and perform more DB interaction 
server side that bar.php. Therefor, we limit our tests to one resource. 

The tests are performed over the Internet, no WiFi, against a server with
no other significant traffic than the test traffic itself. 

The IO model is blocking, non-multiplexing within the concurrent processes.
This is a common concurrency model for threaded applications running on e.g.,
.NET CLR or JVM platforms.


## weblat.go measurements

- Low to high NCONCURRENTMAX with low NREQS (bursts) against static content 
- Low to high NCONCURRENTMAX with high NREQS (sustained traffic) against static 
  content
- Low to high NCONCURRENTMAX with low NREQS (bursts) against dynamic content 
- Low to high NCONCURRENTMAX with high NREQS (sustained traffic) against 
  dynamic content
