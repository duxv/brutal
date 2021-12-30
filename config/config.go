package config

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/duxv/brutal/logging"
)

var validMethods = []string{
	"GET", "HEAD", "POST", "PUT",
	"DELETE", "CONNECT", "OPTIONS",
	"TRACE", "PATCH",
}

type Matcher struct {
	Length      int
	Regex       *regexp.Regexp
	StatusCodes []int
}

type Config struct {
	// the raw target address
	target string
	// the parsed list of every path to try
	wordlist []string
	// the amount of attempts to run
	// at the same time, recommending
	// not setting it to a high number
	threadCount int
	// amount of time to wait
	// (in milliseconds)
	timeout int
	// method type of the request
	method string
	// the matching configuration
	matcher *Matcher
}

func isValidUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func New(target string) (*Config, error) {

	if !isValidUrl(target) {
		return nil, errors.New("invalid URL")
	}

	c := &Config{
		target:  target,
		matcher: &Matcher{},
	}

	return c, nil
}

// Set the timeout to a specific number of milliseconds
func (c *Config) SetTimeout(milliseconds int) {
	c.timeout = milliseconds
}

// Set the number of goroutines running at the same time
func (c *Config) SetThreadCount(count int) {
	c.threadCount = count
}

// Set the method of the request
func (c *Config) SetMethod(method string) {
	has := false

	for _, m := range validMethods {
		if m == method {
			has = true
			break
		}
	}

	if !has {
		logging.Critical("invalid request method %s", method)
	}
	c.method = method
}

// Set the length property of the matcher
func (c *Config) SetMatcherLength(num int) {
	c.matcher.Length = num
}

// Set the regex property of the matcher
func (c *Config) SetMatcherRegex(re *regexp.Regexp) {
	c.matcher.Regex = re
}

func (c *Config) AddMatcherStatusCodes(codes ...int) {
	c.matcher.StatusCodes = append(c.matcher.StatusCodes, codes...)
}

// Add one or multiple words to the wordlist
func (c *Config) AddWord(word ...string) {
	for _, w := range word {
		w = strings.TrimSpace(w)
		if w != "" {
			c.wordlist = append(c.wordlist, w)
		}
	}
}

// Add one or multiple valid codes to the valid
func (c *Config) AddMatcherStatusCodesString(codes ...string) {
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		num, err := strconv.Atoi(code)
		if err != nil {
			logging.Critical("Error while trying to parse status code '%s': %v", code, err)
		}
		c.matcher.StatusCodes = append(c.matcher.StatusCodes, num)
	}
}

func (c Config) Threads() int       { return c.threadCount }
func (c Config) Timeout() int       { return c.timeout }
func (c Config) Wordlist() []string { return c.wordlist }
func (c Config) Target() string     { return c.target }
func (c Config) Method() string     { return c.method }
func (c Config) Matcher() *Matcher  { return c.matcher }
