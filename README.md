# unimport

unimport is a Go static analysis tool to find unnecessary import aliases.

## Installation

    go get -u github.com/alexkohler/unimport

## Usage

Similar to other Go static anaylsis tools (such as golint, go vet) , unimport can be invoked with one or more filenames, directories, or packages named by its import path. Unimport also supports the `...` wildcard. 

    unimport files/directories/packages

Currently, no flag are supported. A `-w` flag may be added in the future to automatically remove aliases where possible. (Similar to [gofmt's -w flag](https://golang.org/cmd/gofmt/))

## Purpose

As noted in Go's [Code Review comments](https://github.com/golang/go/wiki/CodeReviewComments#imports):

> Avoid renaming imports except to avoid a name collision; good package names should not require renaming. 
> In the event of collision, prefer to rename the most local or project-specific import.

This tool will check if any import aliases are truly needed (by ensuring there is a name collision that would exist without the mport alias). Furthermore, unimport will flag any use of the [import dot](https://github.com/golang/go/wiki/CodeReviewComments#import-dot) outside of test files.

## Example

Running `unimports` on the [Go source](https://github.com/golang/go):

```Bash
$ unimport $GOROOT/src/...
//TODO finish
```

Below is... //TODO


```Go
//TODO
```

## TODO

- Unit tests (may require some refactoring to do correctly)
- -w flag to write changes to file where/if possible
- Vim quickfix format?
- Globbing support (e.g. unimport *.go)


## Contributing

Pull requests welcome!

