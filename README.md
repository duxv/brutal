# Brutal

A lightweight, simple to use web fuzzer.


# Usage

<p> Brutal is pretty easy to use. 
  <ol> 
  <li> <b>--debug</b> - for printing more details about the runtime (do not use if you are not trying to develop the program) </li>
  <li> <b>--help</b> - retrieve all the commands and a small description </li>
  <li> <b>--method</b> - change the method of the done requests *must be valid* </li>
  <li> <b>--threads</b> - the amount of requests to be done at the same time;
  it is recommended to not use too many. </li>
    <li> <b>--timeout</b> - the amount of seconds to wait for a request; if timeout is reached, the request is not count as valid </li>
    <li> <b>--match-status</b> - the status code of a request must be part of these or it will be count as invalid; if not specified defaults are being used (statuses must be separated by commas); to completely ignore the status use --match-status "" </li>
    <li> <b>--match-length</b> - the response body must have this amount of characters; default is disabled</li>
    <li> <b>--match-regex</b> - the response body must match this regex</li>
    <li> <b>--wordlist</b> - the file to retrieve the words from</li>
    <li> <b>--wordlist-separator</b> what to separate the words by in the wordlist; by default they are separated by a newline </li>

   </ol>
</p>

### FUZZ keyword

<p>The 'FUZZ' keyword is used to represent the place where to word is going to be <br> 

For instance, if you run `brutal http://localhost:9000/FUZZ`, and the words in the wordlist are 'etc' and 'passwd', the next URLs will be requested: <br>
- `http://localhost:9000/etc`
- `http://localhost:9000/passwd`

<br>

If you'd use the keyword multiple times, it will get replaced each time it has been used.

</p>

### What amount of threads to use?

Sometimes one thread can be faster than four, because of the synchronization price. <br>
If you have 0ms latency, one thread would be just fine, either way you can use more. <br>
It depends on the time it takes to process a request.<br>

## Having a suggestion?

If you have any suggestions, ideas or found any bug, you can join [this Discord server](https://discord.gg/ktEBKceytN).<br>

### Output of the --help flag
```
Usage:
  brutal [flags]

Examples:
brutal http://127.0.0.1/FUZZ

Flags:
  -d, --debug                       print more information about the runtime
  -h, --help                        help for brutal
  -l, --match-length int            length of the response body must be equal to (default -1)
  -r, --match-regex string          response body must match this regex
  -x, --match-status string         http status codes identified as valid (separated by a comma) (default "200,201,202,203,204,205,206,207,208,226")
  -m, --method string               method of the requests to be done (default "GET")
  -p, --threads int                 number of attempts to run at the same time (default 1)
  -t, --timeout int                 time in seconds to wait for a request (default 5)
  -w, --wordlist string             words to use
  -s, --wordlist-separator string   separator of words in the wordlist (default "\n")
```

# Warning

The tool is still in the very early development stage and critical bugs can occurr.