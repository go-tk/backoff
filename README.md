# backoff

[![GoDev](https://pkg.go.dev/badge/golang.org/x/pkgsite.svg)](https://pkg.go.dev/github.com/go-tk/backoff)
[![Build Status](https://travis-ci.org/go-tk/backoff.svg?branch=master)](https://travis-ci.org/github/go-tk/backoff)
[![Coverage Status](https://codecov.io/gh/go-tk/backoff/branch/master/graph/badge.svg)](https://codecov.io/gh/go-tk/backoff)

Exponential backoff algorithm with randomized jitter

## Example

```go
package main

import (
        "context"
        "fmt"
        "log"
        "net/http"
        "time"

        "github.com/go-tk/backoff"
)

func main() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        _ = cancel
        backoff := backoff.New(backoff.Options{
                MinDelay:            100 * time.Millisecond,
                MaxDelay:            100 * time.Second,
                DelayFactor:         2,
                MaxDelayJitter:      1,
                DelayFunc:           backoff.DelayWithContext(ctx),
                MaxNumberOfAttempts: 100,
        })
        req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com/", nil)
        if err != nil {
                log.Fatal(err)
        }
        for {
                resp, err := http.DefaultClient.Do(req)
                if err != nil {
                        if err2 := backoff.Do(); err2 != nil {
                                log.Printf("failed to back off; err=%q", err2)
                                log.Fatal(err)
                        }
                        continue
                }
                resp.Body.Close()
                if resp.StatusCode/100 == 5 {
                        err := fmt.Errorf("http server error; httpStatusCode=%v", resp.StatusCode)
                        if err2 := backoff.Do(); err2 != nil {
                                log.Printf("failed to back off; err=%q", err2)
                                log.Fatal(err)
                        }
                        continue
                }
                fmt.Println(resp.StatusCode)
                return
        }
        // Output:
        // 200
}
```
