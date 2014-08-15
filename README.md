Ocean
=====
A simple lexer for go that supports shell-style quoting, commenting, piping, redirecting, and escaping. Ocean uses a simple state machine to tokenize a complete shell command into individual words and sentences. For example the shell command `ls -l . | grep "file name" > output.txt` is tokenized into the array `["ls", "-l", ".", "|", "grep", "file name", ">", "output.txt"]`. Note that Ocean does more than split the string on each splace. It's mindful of quotes and escaping.


### Examples
```go
str := `one two three`
val, _ := ocean.Tokenize(str)
fmt.Printf("%q", val)
// outputs: ["one", "two", "three"]
```


```go
str := `one "two three"`
val, _ := ocean.Tokenize(str)
fmt.Printf("%q", val)
// outputs: ["one", "two three"]
```


```go
str := `one|two three`
val, _ := ocean.Tokenize(str)
fmt.Printf("%q", val)
// outputs: ["one", "|", "two", "three"]
```


```go
str := `one | "two three" > output.txt`
val, _ := ocean.Tokenize(str)
fmt.Printf("%q", val)
// outputs: ["one", "|", "two three", ">", "output.txt"]
```

### Installing
Use the 'go get' command to download the package.
```bash
go get github.com/headzoo/ocean
```

Then import the package into your project.
```go
import "github.com/headzoo/ocean"
```


### Credits
Ocean is an indirect fork of [go-shlex](http://code.google.com/p/go-shlex/) written by Steven Thurgood. Contributions were make by Jonathan Rudenberg through his fork at [flynn/go-shlex](https://github.com/flynn/go-shlex).
