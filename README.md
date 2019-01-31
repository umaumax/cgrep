# cgrep

color enhanced grep

## how to use
```
git grep --color=always . | cgrep '(".*")' 'yellow' | cgrep '(%[0-9]*[dsf])' 'magenta' | cgrep '([0-9]+)' green
```

## how to install
```
go get -u "github.com/umaumax/cgrep"
```

## color format info
* [mgutz/ansi: Small, fast library to create ANSI colored strings and codes\. \[go, golang\]]( https://github.com/mgutz/ansi )
