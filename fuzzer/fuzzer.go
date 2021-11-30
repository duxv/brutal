package fuzzer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/duxv/brutal/config"
	"github.com/duxv/brutal/logging"
	"github.com/logrusorgru/aurora"
)

var (
	mutex sync.Mutex
)

type Result struct {
	address string
}
type Fuzzer struct {
	matcher  *config.Matcher
	method   string
	target   string
	threads  int
	timeout  int
	wordlist []string
}

func New(conf *config.Config) *Fuzzer {
	f := &Fuzzer{
		method:   conf.Method(),
		threads:  conf.Threads(),
		timeout:  conf.Timeout(),
		wordlist: conf.Wordlist(),
		target:   conf.Target(),
		matcher:  conf.Matcher(),
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
		host := strings.Replace(f.target, "FUZZ", w, -1)
		mutex.Unlock()
		logging.Debug("Trying %s", host)
		res, err := f.request(host)
		if err != nil {
			results <- nil
			continue
		}
		body, _ := ioutil.ReadAll(res.Body)
		bodyStr := string(body)

		// Check if the response meets all of the matcher's requirements
		switch {
		case f.matcher.Length > -1 && len([]rune(bodyStr)) != f.matcher.Length:
			results <- nil
			continue
		case len(f.matcher.StatusCodes) > 0 && !containsInt(f.matcher.StatusCodes, res.StatusCode):
			results <- nil
			continue
		case f.matcher.Regex != nil && !f.matcher.Regex.Match(body):
			results <- nil
			continue
		}

		results <- &Result{
			address: host,
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
	logging.Info("Target URL: %s", strings.Replace(f.target, "FUZZ", "{{.word}}", -1))

	if len(f.matcher.StatusCodes) > 0 {
		logging.Info("Positive codes: %v", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(f.matcher.StatusCodes)), ", "), "[]"))
	}
	if f.matcher.Length != -1 {
		logging.Info("Exact length: %d", f.matcher.Length)
	}
	if f.matcher.Regex != nil {
		logging.Info("Response body match %q", f.matcher.Regex.String())
	}
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
		if r != nil {
			fmt.Printf("[%v] %s\n", aurora.Blue("MATCH"), r.address)
		}
	}
}
