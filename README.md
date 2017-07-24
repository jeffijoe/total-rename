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

# Author

Jeff Hansen - [@Jeffijoe](https://twitter.com/Jeffijoe)

> This is my very first Go project, and it is meant as a learning experience
> for trying out the Go language while building something useful (to me, at least).
>
> Simply put, I want to see what the fuss is all about.
> When I started writing this I had written exactly 0 lines of Go code.
>
> This is my **Hello World**.


  [Releases]: https://github.com/jeffijoe/total-rename/releases
  [screenshot]: http://i.imgur.com/3NaGKzT.png
