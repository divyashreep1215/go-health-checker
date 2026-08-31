package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

func checkURL(ctx context.Context, url string, client *http.Client) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{URL: url, Err: err}
	}

	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return Result{URL: url, Duration: duration, Err: err}
	}
	defer resp.Body.Close()

	return Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}

func worker(ctx context.Context, jobs <-chan string, results chan<- Result, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		results <- checkURL(ctx, url, client)
	}
}

func main() {
	targets := []string{
		"https://golang.org",
		"https://github.com",
		"https://google.com",
		"https://httpbin.org/get",
	}

	numJobs := len(targets)
	numWorkers := 3

	jobs := make(chan string, numJobs)
	results := make(chan Result, numJobs)

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(ctx, jobs, results, client, &wg)
	}

	for _, target := range targets {
		jobs <- target
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Printf("%-38s | %-8s | %-12s | %s\n", "URL", "STATUS", "LATENCY", "ERROR")
	fmt.Println("----------------------------------------------------------------------------------")

	for res := range results {
		if res.Err != nil {
			fmt.Printf("%-38s | %-8s | %-12s | %v\n", res.URL, "FAIL", res.Duration.Truncate(time.Millisecond), res.Err)
		} else {
			fmt.Printf("%-38s | %-8d | %-12s | %s\n", res.URL, res.StatusCode, res.Duration.Truncate(time.Millisecond), "OK")
		}
	}
}
