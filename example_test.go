package try_test

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/go-tk/try"
)

func ExampleDo() {
	// Wait for file creation.
	ok, err := Do(context.Background(), func() (bool, error) {
		if _, err := os.Stat("foo.txt"); err != nil {
			if os.IsNotExist(err) {
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
	fmt.Println(ok)
	// Output:
	// false
}
