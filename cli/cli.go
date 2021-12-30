package cli

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/duxv/brutal/config"
	"github.com/duxv/brutal/fuzzer"
	"github.com/duxv/brutal/logging"
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
	// time to wait for a request in milliseconds
	timeout int = 5000
	// the method of the web requests
	method string = "GET"

	// the status codes count as valid
	validCodes string = "200,201,202,203,204,205,206,207,208,226"
	// the length of the response body
	responseLength int = -1
	// the raw regexp match expression
	rawRegex  string
	quickList string
)

var cmd = &cobra.Command{
	Use:     "brutal",
	Example: "brutal http://127.0.0.1/FUZZ",
	Short:   "A bloatless url fuzzer",
	Run: func(cmd *cobra.Command, args []string) {
		// command line errors
		switch {
		case len(args) == 0:
			logging.Critical("Must provide a target URL")
		case wordlistPath == "" && quickList == "":
			logging.Critical("No wordlist provided")
		case threads < 1:
			logging.Critical("Threads cannot be below 0, at least 1.")
		case !strings.Contains(args[0], "FUZZ"):
			logging.Critical(`
			Did not specify keyword FUZZ in the URL
			Example: http://localhost:9000/FUZZ
			Brutal will replace the keyword fuzz every time with a word from wordlist.
			`)
		}

		// create a new config instance from args[0]
		conf, err := config.New(args[0])
		if err != nil {
			logging.Critical("URL parse error: %v", err)
		}
		if wordlistPath != "" {
			wordlistContent, err := ioutil.ReadFile(wordlistPath)
			if err != nil {
				logging.Critical("error trying to open wordlist: '%v'", err)
			}
			wlsContent := string(wordlistContent)
			conf.AddWord(strings.Split(wlsContent, wordlistSeparator)...)
		}
		conf.SetThreadCount(threads)
		conf.SetTimeout(timeout)
		conf.SetMethod(method)
		conf.AddWord(strings.Split(quickList, ",")...)
		conf.AddMatcherStatusCodesString(strings.Split(validCodes, ",")...)
		conf.SetMatcherLength(responseLength)
		if rawRegex != "" {
			compiledRegex, err := regexp.Compile(rawRegex)
			if err != nil {
				logging.Critical("Invalid regex to match: %q", rawRegex)
			}
			conf.SetMatcherRegex(compiledRegex)
		}
		fuzz := fuzzer.New(conf)
		fuzz.Run()
	},
}

func init() {
	cmd.PersistentFlags().IntVarP(&threads, "threads", "p", threads, "number of attempts to run at the same time")
	cmd.PersistentFlags().StringVarP(&method, "method", "m", method, "method of the requests to be done")
	cmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", timeout, "time in milliseconds to wait for a request")
	cmd.PersistentFlags().StringVarP(&wordlistPath, "wordlist", "w", wordlistPath, "words to use")
	cmd.PersistentFlags().StringVarP(&wordlistSeparator, "wordlist-separator", "s", wordlistSeparator, "separator of words in the wordlist")
	cmd.PersistentFlags().BoolVarP(&logging.DebugEnable, "debug", "d", logging.DebugEnable, "print more information about the runtime")
	// cmd.PersistentFlags().StringVarP(&configPath, "config", "c", configPath, "path of the config file")
	cmd.PersistentFlags().StringVarP(&validCodes, "match-status", "x", validCodes, "http status codes identified as valid (separated by a comma)")
	cmd.PersistentFlags().IntVarP(&responseLength, "match-length", "l", responseLength, "length of the response body must be equal to")
	cmd.PersistentFlags().StringVarP(&rawRegex, "match-regex", "r", rawRegex, "response body must match this regex")
	cmd.PersistentFlags().StringVarP(&quickList, "quick-list", "q", quickList, "use a wordlist from the command line arguments (separated by a comma)")
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		logging.Critical("%v", err)
	}
}
