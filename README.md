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

Difference between requests to different resources is expected. foo.php might 
execute more code and perform more DB interaction server side that bar.php.
Here, we assume that we're already 

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

The tests are performed over a wired network.

The IO model is blocking, non-multiplexing within the concurrent processes.
This is a common concurrency model for threaded applications running on e.g.,
.NET CLR or JVM platforms.


## weblat.go measurements

- High NCONCURRENTMAX with low NREQS (bursts) against static content 
- High NCONCURRENTMAX with high NREQS (sustained traffic) against static 
  content
- High NCONCURRENTMAX with low NREQS (bursts) against dynamic content 
- High NCONCURRENTMAX with high NREQS (sustained traffic) against dynamic 
  content
