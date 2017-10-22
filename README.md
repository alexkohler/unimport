# unimport

unimport is a Go static analysis tool to find unnecessary import aliases.

## Installation

    go get -u github.com/alexkohler/unimport

## Usage

Similar to other Go static anaylsis tools (such as golint, go vet) , unimport can be invoked with one or more filenames, directories, or packages named by its import path. Unimport also supports the `...` wildcard. 

    unimport files/directories/packages

Currently, no flag are supported. A `-w` flag may be added in the future to automatically remove aliases where possible (Similar to [gofmt's -w flag](https://golang.org/cmd/gofmt/)).

## Purpose

As noted in Go's [Code Review comments](https://github.com/golang/go/wiki/CodeReviewComments#imports):

> Avoid renaming imports except to avoid a name collision; good package names should not require renaming. 
> In the event of collision, prefer to rename the most local or project-specific import.

This tool will check if any import aliases are truly needed (by ensuring there is a name collision that would exist without the mport alias). This tool _will_ ignore import paths containing dashes and dots, as these are generally useful aliases while importing a specific revision. For example, in [gometalinter](https://github.com/alecthomas/gometalinter), there are some imports like `kingpin "gopkg.in/alecthomas/kingpin.v3-unstable"`. This is a reasonable import alias and will not be flagged.

## Example

Running `unimports` on the [Go source](https://github.com/golang/go):

```
$ unimport $GOROOT/src/...
cmd/go/pkg.go:18 unnecessary import alias pathpkg
go/build/build.go:19 unnecessary import alias pathpkg
go/internal/gcimporter/gcimporter.go:23 unnecessary import alias exact
os/pipe_test.go:14 unnecessary import alias osexec
os/os_windows_test.go:10 unnecessary import alias osexec
```

Below are some of the arguably unneeded import aliases it found:


```Go

// go/build/build.go
import (                                                                                       
    "bytes"                                                                                    
    "errors"                                                                                   
    "fmt"                                                                                      
    "go/ast"                                                                                   
    "go/doc"                                                                                   
    "go/parser"                                                                                
    "go/token"                                                                                 
    "io"                                                                                       
    "io/ioutil"                                                                                
    "log"                                                                                      
    "os"                                                                                       
    pathpkg "path"                                                                             
    "path/filepath"                                                                            
    "runtime"                                                                                  
    "sort"                                                                                     
    "strconv"                                                                                  
    "strings"                                                                                  
    "unicode"                                                                                  
    "unicode/utf8"                                                                             
) 

// go/internal/gcimporter/gcimporter.go
import (                                                                                       
    "bufio"                                                                                    
    "errors"                                                                                   
    "fmt"                                                                                      
    "go/build"                                                                                 
    "go/token"                                                                                 
    "io"                                                                                       
    "io/ioutil"                                                                                
    "os"                                                                                       
    "path/filepath"                                                                            
    "sort"                                                                                     
    "strconv"                                                                                  
    "strings"                                                                                  
    "text/scanner"                                                                             
                                                                                               
    exact "go/constant"                                                                        
    "go/types"                                                                                 
)


// os/pipe_test.go.go
import (                                                                                       
    "fmt"                                                                                      
    "internal/testenv"                                                                         
    "os"                                                                                       
    osexec "os/exec"                                                                           
    "os/signal"                                                                                
    "syscall"                                                                                  
    "testing"                                                                                  
)
```


## TODO

- Unit tests
- Flagging of packages that contain an uppercase letter or underscore
- -w flag to write changes to file where/if possible
- Globbing support (e.g. unimport *.go)


## Contributing

Pull requests welcome!

