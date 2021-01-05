package backoff_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-tk/backoff"
)

func ExampleBackoff() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = cancel
	b := backoff.New(backoff.Options{
		MinDelay:            100 * time.Millisecond,        // default
		MaxDelay:            100 * time.Second,             // default
		DelayFactor:         2,                             // default
		MaxDelayJitter:      1,                             // default
		DelayFunc:           backoff.DelayWithContext(ctx), // with respect to ctx
		MaxNumberOfAttempts: 100,                           // default
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
