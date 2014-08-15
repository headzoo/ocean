Ocean
=====
A simple lexer for go that supports shell-style quoting, commenting, piping, redirecting, and escaping.


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
Ocean is an indirect fork of [go-shlex](http://code.google.com/p/go-shlex/) written by Steven Thurgood. Contributions were make by Jonathan Rudenberg through his fork [flynn/go-shlex](https://github.com/flynn/go-shlex).
