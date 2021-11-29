# Brutal

A lightweight, simple to use web fuzzer.


```
Usage:
  brutal [flags]

Examples:
brutal http://127.0.0.1

Flags:
  -d, --debug                       print more information about the runtime
  -h, --help                        help for brutal
  -m, --method string               method of the requests to be done (default "GET")
  -p, --threads int                 number of attempts to run at the same time (default 1)
  -t, --timeout int                 time in seconds to wait for a request (default 5)
  -v, --valid-codes string          http status codes identified as valid (separated by a comma)
  -w, --wordlist string             words to use
  -s, --wordlist-separator string   separator of words in the wordlist (default "\n")
```

# Warning

The tool is still in the very early development stage and critical bugs can occurr.