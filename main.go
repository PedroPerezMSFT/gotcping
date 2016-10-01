package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/montanaflynn/stats"
)

func main() {

	hostPtr := flag.String("host", "bing.com", "Host or IP address to test")
	portPtr := flag.Int("port", 80, "Port number to query")
	countPtr := flag.Int("count", 10, "Number of requests to send")
	timeoutPtr := flag.Int("timeout", 5, "Timeout for each request, in seconds")

	flag.Parse()

	//args := flag.Args()
	host := *hostPtr
	port := *portPtr
	count := *countPtr
	timeout := *timeoutPtr

	_, err := net.LookupIP(host)
	if err != nil {
		fmt.Println("error: unknown host")
		os.Exit(2)
	}

	ping(host, port, count, timeout)

}

func ping(host string, port int, count int, timeout int) {
	successfulProbes := 0
	i := 1
	timeTotal := time.Duration(0)
	var responseTimes []float64

	addr := fmt.Sprintf("%s:%d", host, port)

	for i = 1; count >= i; i++ {
		timeStart := time.Now()
		_, err := net.DialTimeout("tcp", addr, time.Second*time.Duration(timeout))
		responseTime := time.Since(timeStart)
		if err != nil {
			fmt.Println(fmt.Sprintf("%s port %d closed.", host, port))
		} else {
			fmt.Println(fmt.Sprintf("Connected to %s:%d, RTT=%.2fms", host, port, float32(responseTime)/1e6))
			timeTotal += responseTime
			successfulProbes++
		}

		time.Sleep(1e9)
	}

	// Let's calculate and spill some results
	// 1. Average response time
	timeAverage := time.Duration(int64(timeTotal) / int64(successfulProbes))

	// 2. Min and Max response times
	var biggest float64

	smallest := float64(1000000000)

	for _, v := range responseTimes {

		if v > biggest {
			biggest = v
		}

		if v < smallest {
			smallest = v
		}

	}

	// 3. Median response time
	median, _ := stats.Median(responseTimes)

	// 4. Percentile
	percentile90, _ := stats.Percentile(responseTimes, float64(90))
	percentile75, _ := stats.Percentile(responseTimes, float64(75))
	percentile50, _ := stats.Percentile(responseTimes, float64(50))
	percentile25, _ := stats.Percentile(responseTimes, float64(25))

	fmt.Println("\nProbes sent:", i-1, "\nSuccessful responses:", successfulProbes, "\n% of requests failed:", float64(100-(successfulProbes*100)/(i-1)), "\nMin response time:", time.Duration(smallest), "\nAverage response time:", timeAverage, "\nMedian response time:", time.Duration(median), "\nMax response time:", time.Duration(biggest))

	fmt.Println("\n90% of requests were faster than:", time.Duration(percentile90), "\n75% of requests were faster than:", time.Duration(percentile75), "\n50% of requests were faster than:", time.Duration(percentile50), "\n25% of requests were faster than:", time.Duration(percentile25))

}
