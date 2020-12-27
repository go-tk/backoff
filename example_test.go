package try_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	. "github.com/go-tk/try"
)

func ExampleDo() {
	var resp *http.Response
	ok, err := Do(context.Background(), func() (bool, error) {
		var err error
		resp, err = http.Get("http://example.com")
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Temporary() {
				return false, nil // retry
			}
			return false, err // error
		}
		return true, nil // succeed
	}, Options{
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
