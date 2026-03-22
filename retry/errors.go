package retry

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"syscall"
)

// HTTPError signals a retryable HTTP status code. Callers wrap this error
// after inspecting the response status so the retry loop can classify it.
type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d", e.StatusCode)
}

// IsRetryable returns true for transient network errors and HTTP 429 / 5xx.
// It does NOT retry 4xx (except 429), JSON decode errors, or validation errors.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for HTTPError (429 or 5xx).
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == 429 || httpErr.StatusCode >= 500
	}

	// Transient network errors (timeouts, temporary failures).
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	// URL errors wrapping network failures.
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return true
	}

	// Specific syscall errors indicating connection problems.
	if errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.ECONNREFUSED) ||
		errors.Is(err, syscall.EPIPE) {
		return true
	}

	return false
}
