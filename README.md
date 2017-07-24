<center>

# total-rename

Utility to rename occurences of a string in files in the correct casing — content _and_ path!

![Total Rename][screenshot]

</center>

# Installation

Check the [Releases] section for binary downloads for your platform.

Alternatively, if you have Go installed and configured, you can use

```
go get github.com/jeffijoe/total-rename
```

# Usage

```
total-rename - case-preserving renaming utility
Copyright © Jeff Hansen 2017 to present. All rights reserved.

OPTIONS:
    Options must be specified before arguments.

    --dry         If set, won't rename anything
    --force       Replaces all occurences without asking
    --help        Shows this help text

ARGUMENTS:

    <pattern>  Search pattern (glob). Relative to working
               directory unless rooted (absolute path).
    <find>     The string to find. If multiple words,
               please use camelCase.
    <replace>  The string to replace occurences with.
               If multiple words, please use camelCase.

EXAMPLE:

    total-rename "**/*.txt" "awesome" "excellent"

    Rename all occurences of "awesome" to "excellent" in
    all .txt files (and folders) recursively from the
    current working directory:

EXAMPLE:

    total-rename --force "/Users/jeff/projects/my-app/src/**/*.*" "awesome" "excellent"

    Like the first example, but from an absolute path, and match all
    file extensions and don't ask for confirmation for each occurence.
```

# How it works

`total-rename` will scan every file matched by the pattern you specify, and look for every occurence 
of the search string in every casing format. This works by taking the search string and converting it to
different casings to search for. **The generated casings may be inaccurate depending on the input string**, and
it would seem **the most accurate casings are generated when the input is `camelCased`.** _This also applies
to the replacement string._

After having collected every occurence of the string within every file's content and path, you have the option to
review every change in an interactive way. **Nothing is replaced until the interactive yes-no session is done.**
If you don't want to review every change, you can pass the `--force` flag.

# About

This was my very first Go project, and it was meant as a learning experience
for trying out the Go language while building something useful that I needed.

Simply put, I wanted to see what the fuss was all about. When I started writing 
this I had written exactly 0 lines of Go code.

Things I wanted to cover in this project in order to learn Go was:

* Basic types
* String manipulation
* File I/O
* Goroutines
* Channels
* Splitting work to run in parallel with goroutines and channels, avoiding deadlocks
* Testing
* Using external dependencies
* Manipulating the terminal with colors and moving the cursor <small>(holy shit that was cumbersome)</small>
* Accepting user input <small>(`fmt.Scanln()` does not do what I think it does...)</small>
* Cross-compilation

**Disclaimer:** the following paragraphs describe my experience using Go, and it's not all good. **I am not saying "Go sucks!", I am just pointing out my personal disappointments as a Node/.NET developer**.

I have to say, after having written this project in Go, I have a renewed appreciation for JavaScript and everything you get for free, including (but not limited to) filtering/mapping arrays, arrow functions, Promises, and the wealth of small modules available on npm. Every time I had to declare a `result` array, then a `for ... range` loop that `result = append(result, ...)` just to map from 1 thing to another, a tiny piece of my soul died.

I'm using macOS as my development machine, but I want to target Windows and Linux as well, so it's awesome that Go supports cross compilation! :+1: `goreleaser` makes this even better by managing everything related to building for different OS'es and architectures.

Tests run super fast, which is nice! :+1: However writing the tests was pretty weird, having to pass in `t` to `assert.X(t, ...)` felt awkward. And the output from `go test` is not very human friendly; thankfully [richgo] made it a little more readable by coloring the output.

Using `go get` to manage packages is nice, but pulling the latest `master` does not seem very production-friendly to me. Sure, `master` is supposed to be stable at all costs, but humans are not perfect and accidental breaking changes slip in. It's awesome that `go get` is built in, but I think version locking is important.

I was using VS Code for writing Go, and the Go extension is awesome; auto imports is nice when it works (it won't work if your code does not compile 😞).

I've been used to using 2-space indentation for years - having to use 8 tabs really grinds my gears. I'm all for having a code standard enforced by the official tooling, but 8 tabs is crazy if you ask me — so much screen real-estate goes to waste. Please, at least make it 4... _spaces!_ 🙃

**With all this being said,** for building performance-critical system components, I would definitely consider using Go! While the syntax is lacking, such as arrow funcs, more type inference, and basic functional programming things like array `map`, `filter` and `reduce`, the performance is great and goroutines + channels + `sync.WaitGroup` is awesome!

# Author

Jeff Hansen - [@Jeffijoe](https://twitter.com/Jeffijoe)

  [Releases]: https://github.com/jeffijoe/total-rename/releases
  [screenshot]: http://i.imgur.com/3NaGKzT.png
  [richgo]: https://github.com/kyoh86/richgo
