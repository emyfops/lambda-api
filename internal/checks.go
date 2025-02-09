package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func HTTPGetCheck(url string) func(ctx context.Context) error {
	client := http.Client{
		// never follow redirects
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return func(ctx context.Context) error {
		// Retrieve and set the timeout from the health library
		deadline, _ := ctx.Deadline()
		client.Timeout = time.Since(deadline)

		resp, err := client.Get(url)
		if err != nil {
			return err
		}

		if resp.StatusCode < 200 && resp.StatusCode >= 300 {
			return fmt.Errorf("returned status %d", resp.StatusCode)
		}
		return nil
	}
}
