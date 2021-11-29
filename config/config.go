package config

import (
	"brutal/logging"
	"errors"
	"net/url"
	"strings"
)

var defaultCodes = []int{
	200, 201, 202, 203, 204,
	205, 206, 207, 208, 226,
}

type Config struct {
	// the raw target address
	target string
	// will count as successful only
	// if the status code is part of
	// one of those integers
	validCodes []int
	// the parsed list of every path to try
	wordlist []string
	// the amount of attempts to run
	// at the same time, recommending
	// not setting it to a high number
	threadCount int
	// amount of time to wait
	// (in seconds)
	timeout int
}

func isValidUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func New(target string) (*Config, error) {
	if !strings.HasSuffix(target, "/") {
		target += "/"
	}

	if !isValidUrl(target) {
		return nil, errors.New("invalid URL")
	}

	c := &Config{target: target}

	return c, nil
}

// Set the timeout to a specific number of seconds
func (c *Config) SetTimeout(seconds int) {
	c.timeout = seconds
}

// Set the number of goroutines running at the same time
func (c *Config) SetThreadCount(count int) {
	c.threadCount = count
}

// Add one or multiple words to the wordlist
func (c *Config) AddWord(word ...string) {
	for _, w := range word {
		if w != "" && w != "\t" && w != "\r" {
			c.wordlist = append(c.wordlist, w)
		}
	}
}

// Add one or multiple valid codes to the valid codes
func (c *Config) AddValidCode(code ...int) {
	c.validCodes = append(c.validCodes, code...)
}

// Change all the settings to the default ones
func (c *Config) ResetDefault() {
	c.AddValidCode(defaultCodes...)
	logging.Debug("Added all default status codes")
}

func (c *Config) Threads() int       { return c.threadCount }
func (c *Config) Timeout() int       { return c.timeout }
func (c *Config) Wordlist() []string { return c.wordlist }
func (c *Config) ValidCodes() []int  { return c.validCodes }
func (c *Config) Target() string     { return c.target }
