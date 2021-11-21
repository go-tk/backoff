# backoff

[![GoDev](https://pkg.go.dev/badge/golang.org/x/pkgsite.svg)](https://pkg.go.dev/github.com/go-tk/backoff
) [![Workflow Status](https://github.com/go-tk/optional/actions/workflows/default.yaml/badge.svg)](https://github.com/go-tk/optional/actions
) [![Coverage Status](https://codecov.io/gh/go-tk/backoff/branch/master/graph/badge.svg)](https://codecov.io/gh/go-tk/backoff
)

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
        "github.com/go-tk/optional"
)

func main() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        _ = cancel
        b := backoff.New(backoff.Options{
                MinDelay:            optional.MakeDuration(100 * time.Millisecond), // default
                MaxDelay:            optional.MakeDuration(100 * time.Second),      // default
                DelayFactor:         optional.MakeFloat64(2),                       // default
                MaxDelayJitter:      optional.MakeFloat64(1),                       // default
                DelayFunc:           backoff.DelayWithContext(ctx),                 // with respect to ctx
                MaxNumberOfAttempts: optional.MakeInt(100),                         // default
        })
        req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com/", nil)
        if err != nil {
                log.Fatal(err)
        }
        for {
                resp, err := http.DefaultClient.Do(req)
                if err != nil {
                        if err2 := b.Do(); err2 != nil { // delay
                                log.Printf("failed to back off; err=%q", err2)
                                log.Fatal(err)
                        }
                        continue // retry
                }
                resp.Body.Close()
                if resp.StatusCode/100 == 5 {
                        err := fmt.Errorf("http server error; httpStatusCode=%v", resp.StatusCode)
                        if err2 := b.Do(); err2 != nil { // delay
                                log.Printf("failed to back off; err=%q", err2)
                                log.Fatal(err)
                        }
                        continue // retry
                }
                fmt.Println(resp.StatusCode)
                return
        }
        // Output:
        // 200
}
```
