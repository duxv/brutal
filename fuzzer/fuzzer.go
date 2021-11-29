package fuzzer

import (
	"brutal/config"
	"brutal/logging"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/logrusorgru/aurora"
)

var (
	mutex sync.Mutex
)

type Result struct {
	status int
	path   string
}
type Fuzzer struct {
	method     string
	target     string
	threads    int
	timeout    int
	wordlist   []string
	validCodes []int
}

func New(conf *config.Config) *Fuzzer {
	f := &Fuzzer{
		method:     conf.Method(),
		threads:    conf.Threads(),
		timeout:    conf.Timeout(),
		wordlist:   conf.Wordlist(),
		validCodes: conf.ValidCodes(),
		target:     conf.Target(),
	}
	return f
}

func containsInt(arr []int, e int) bool {
	for _, element := range arr {
		if element == e {
			return true
		}
	}
	return false
}

func (f *Fuzzer) worker(words <-chan string, results chan<- *Result) {
	for w := range words {
		mutex.Lock()
		host := f.target + w
		mutex.Unlock()
		req, err := f.request(host)
		if err != nil {
			results <- nil
			continue
		}
		results <- &Result{
			status: req.StatusCode,
			path:   w,
		}
	}
}

func (f *Fuzzer) request(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(f.timeout),
	}
	req, err := http.NewRequest(f.method, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	return resp, err
}

func (f *Fuzzer) Run() {
	logging.Info("Fuzzing started")
	logging.Info("Target URL: %s", f.target)
	logging.Info("Positive codes: %v", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(f.validCodes)), ", "), "[]"))
	logging.Info("All results will be printed below. If nothing is printed, means nothing found.")
	words := make(chan string)
	results := make(chan *Result)
	go func() {
		defer close(words) // close words channel as soon as it reads every wordlist
		for _, word := range f.wordlist {
			words <- word
		}
	}()

	go func() {
		defer close(results) // closes results as soon as all workers are done (after words channel is closed)
		wg := sync.WaitGroup{}
		wg.Add(f.threads)
		for i := 0; i < f.threads; i++ {
			go func() {
				defer wg.Done()
				f.worker(words, results) // these will be done as soon as processes all words and channel is closed
			}()
		}
		wg.Wait()
	}()

	for r := range results {
		if r != nil && containsInt(f.validCodes, r.status) {
			fmt.Printf("%v%d%v %s\n", aurora.Yellow("["), aurora.Blue(r.status), aurora.Yellow("]"), f.target+r.path)
		}
	}
}
