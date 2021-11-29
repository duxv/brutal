package cli

import (
	"brutal/config"
	"brutal/fuzzer"
	"brutal/logging"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// the numbers of threads to use
	threads int = 1
	// the path of the wordlist to use
	wordlistPath string
	// the separator to identify words in
	// the wordlist by
	wordlistSeparator string = "\n"
	// // the config file to use
	// // if it is empty, the program will
	// // use the default settings
	// configPath string
	// time to wait for a request in seconds
	timeout    int = 5
	validCodes string
	// the method of the web requests
	method string = "GET"
)

var cmd = &cobra.Command{
	Use:     "brutal",
	Example: "brutal http://127.0.0.1",
	Short:   "A bloatless url fuzzer",
	Run: func(cmd *cobra.Command, args []string) {
		// command line errors
		switch {
		case len(args) == 0:
			logging.Critical("Must provide a target URL")
		case wordlistPath == "":
			logging.Critical("No wordlist provided")
		case threads == 0:
			logging.Critical("Threads cannot be 0, at least 1.")
		}

		// create a new config instance from args[0]
		conf, err := config.New(args[0])
		if err != nil {
			logging.Critical("URL parse error: %v", err)
		}
		wordlistContent, err := ioutil.ReadFile(wordlistPath)
		if err != nil {
			logging.Critical("error trying to open wordlist: '%v'", err)
		}
		wlsContent := string(wordlistContent)
		conf.AddWord(strings.Split(wlsContent, wordlistSeparator)...)
		conf.SetThreadCount(threads)
		conf.SetTimeout(timeout)
		conf.SetMethod(method)
		if validCodes != "" {
			statusCodes := parseArrayInt(validCodes)
			conf.AddValidCode(statusCodes...)
		} else {
			conf.ResetDefault()
		}
		fuzz := fuzzer.New(conf)
		fuzz.Run()
	},
}

func parseArrayInt(v string) []int {
	arr := strings.Split(v, ",")
	intArr := []int{}

	for _, s := range arr {
		s = strings.TrimSpace(s)
		num, err := strconv.Atoi(s)
		if err != nil {
			logging.Critical("invalid number: %s", s)
		}
		intArr = append(intArr, num)
	}
	return intArr
}

func init() {
	cmd.PersistentFlags().IntVarP(&threads, "threads", "p", threads, "number of attempts to run at the same time")
	cmd.PersistentFlags().StringVarP(&method, "method", "m", method, "method of the requests to be done")
	cmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", timeout, "time in seconds to wait for a request")
	cmd.PersistentFlags().StringVarP(&wordlistPath, "wordlist", "w", wordlistPath, "words to use")
	cmd.PersistentFlags().StringVarP(&wordlistSeparator, "wordlist-separator", "s", wordlistSeparator, "separator of words in the wordlist")
	cmd.PersistentFlags().BoolVarP(&logging.DebugEnable, "debug", "d", logging.DebugEnable, "print more information about the runtime")
	// cmd.PersistentFlags().StringVarP(&configPath, "config", "c", configPath, "path of the config file")
	cmd.PersistentFlags().StringVarP(&validCodes, "valid-codes", "v", validCodes, "http status codes identified as valid (separated by a comma)")
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		logging.Critical("%v", err)
	}
}
