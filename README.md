# try

[![Build Status](https://travis-ci.org/go-tk/try.svg?branch=master)](https://travis-ci.org/github/go-tk/try) [![Coverage Status](https://codecov.io/gh/go-tk/try/branch/master/graph/badge.svg)](https://codecov.io/gh/go-tk/try)

Exponential backoff algorithm with jitter

## Documentation

See https://godoc.org/github.com/go-tk/try for details.

## Example

```go
package main

import (
        "context"
        "fmt"
        "net/http"
        "strings"
        "time"

        "github.com/go-tk/try"
)

func main() {
        var resp *http.Response
        ok, err := try.Do(context.Background(), func() (bool, error) {
                var err error
                resp, err = http.Get("http://example.com")
                if err != nil {
                        if err, ok := err.(net.Error); ok && err.Temporary() {
                                return false, nil // retry
                        }
                        return false, err // error
                }
                return true, nil // succeed
        }, try.Options{
                MinBackoff:          100 * time.Millisecond,
                MaxBackoff:          5 * time.Second,
                BackoffFactor:       2,
                MaxNumberOfAttempts: 3,
        })
        if err != nil {
                panic(err)
        }
        if !ok {
                // MaxNumberOfAttempts reached
                return
        }
        resp.Body.Close()
        fmt.Println(resp.StatusCode)
        // Output:
        // 200
}
```
