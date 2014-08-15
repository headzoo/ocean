Ocean
=====
Ocean is a simple lexer for go that supports shell-style quoting, commenting, piping, redirecting, and escaping.


### Examples
```
one two three
```
["one", "two", "three"]

```
one "two three"
```

["one", "two three"]

```
one | two three
```
["one", "|", "two", "three"]

```
one | "two three" > output.txt
```
["one", "|", "two three", ">", "output.txt"]


### Credits
Ocean is an indirect fork of [go-shlex](http://code.google.com/p/go-shlex/) written by Steven Thurgood.
