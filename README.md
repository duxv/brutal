# Brutal

A lightweight, very fast & simple to use web fuzzer.

# Installation

### Install Go

Debian / Ubuntu
```
sudo apt update && sudo apt upgrade && sudo apt install golang-go
```

Arch
```
sudo pacman -Sy go
```


Windows<br>
[Click here](https://go.dev/dl/)

### Install Brutal
```
go install github.com/duxv/brutal@latest
```

# Usage

<p> Brutal is pretty easy to use. </p>


  Command                     |    Description
  ---                         |    ---
  **--debug**                 |    print more details about the runtime
  **--help**                  |    retrieve all commands and a small description for each
  **--method**                |    change the method of the requests
  **--quick-list**            |    use a command line string separated by commas as a wordlist
  **--threads**               |    the amount of requests to be done at the same time
  **--timeout**               |    the amount of milliseconds to wait for a request
  **--match-status**          |    all positive status codes, separated by a comma
  **--match-length**          |    the response body must have this amount of characters
  **--match-regex**           |    the response body must match this regex
  **--wordlist**              |    the file to retrieve the words from
  **--wordlist-separator**    |    what to separate the words by in the wordlist

### FUZZ keyword

<p>The 'FUZZ' keyword is used to represent the place where to word is going to be <br> 

For instance, if you run `brutal http://localhost:9000/FUZZ`, and the words in the wordlist are 'etc' and 'passwd', the next URLs will be requested: <br>
- `http://localhost:9000/etc`
- `http://localhost:9000/passwd`

<br>

If you'd use the keyword multiple times, it will get replaced each time.

</p>

### What amount of threads to use?

Sometimes one thread can be faster than four, because of the synchronization price. <br>
If you have 0ms latency, one thread would be just fine, either way you can use more. <br>
It depends on the time it takes to process a request.<br>

### Which are the valid methods?


<ul> 
  <li>GET</li>
  <li>HEAD</li>
  <li>POST</li>
  <li>PUT</li>
  <li>DELETE</li>
  <li>CONNECT</li>
  <li>OPTIONS</li>
  <li>TRACE</li>
  <li>PATCH</li>
</ul>

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
  -q, --quick-list string           use a wordlist from the command line arguments (separated by a comma)
  -p, --threads int                 number of attempts to run at the same time (default 1)
  -t, --timeout int                 time in milliseconds to wait for a request (default 5000)
  -w, --wordlist string             words to use
  -s, --wordlist-separator string   separator of words in the wordlist (default "\n")
```

# Warning

The tool is still in the very early development stage and critical bugs can occurr.
