package healthz

import (
	"fmt"
	"github.com/heptiolabs/healthcheck"
	"net/http"
	"time"
)

func HTTPGetCheck(url string, timeout time.Duration) healthcheck.Check {
	client := http.Client{
		Timeout: timeout,
		// never follow redirects
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return func() error {
		resp, err := client.Get(url)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode < 200 && resp.StatusCode >= 300 {
			return fmt.Errorf("returned status %d", resp.StatusCode)
		}
		return nil
	}
}
