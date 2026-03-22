package retry

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"io"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestDo_SuccessFirstAttempt(t *testing.T) {
	calls := 0
	err := Do(func() error {
		calls++
		return nil
	}, WithMaxAttempts(3), WithInitialDelay(time.Millisecond))

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetryOnTransientError(t *testing.T) {
	calls := 0
	err := Do(func() error {
		calls++
		if calls < 3 {
			return &net.OpError{Op: "read", Err: fmt.Errorf("connection reset")}
		}
		return nil
	}, WithMaxAttempts(3), WithInitialDelay(time.Millisecond))

	if err != nil {
		t.Fatalf("expected nil error after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_NoRetryOnPermanentError(t *testing.T) {
	calls := 0
	permanent := fmt.Errorf("validation error: bad input")

	err := Do(func() error {
		calls++
		return permanent
	}, WithMaxAttempts(3), WithInitialDelay(time.Millisecond))

	if err != permanent {
		t.Fatalf("expected permanent error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call (no retries for permanent error), got %d", calls)
	}
}

func TestDo_ExhaustsAllAttempts(t *testing.T) {
	calls := 0
	err := Do(func() error {
		calls++
		return &HTTPError{StatusCode: 503}
	}, WithMaxAttempts(3), WithInitialDelay(time.Millisecond))

	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != 503 {
		t.Fatalf("expected HTTPError 503, got %v", err)
	}
}

func TestDoWithResponse_Success(t *testing.T) {
	resp, err := DoWithResponse(func() (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("ok")),
		}, nil
	}, WithMaxAttempts(3), WithInitialDelay(time.Millisecond))

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestDoWithResponse_Retries429(t *testing.T) {
	calls := 0
	resp, err := DoWithResponse(func() (*http.Response, error) {
		calls++
		if calls < 2 {
			return &http.Response{
				StatusCode: 429,
				Body:       io.NopCloser(strings.NewReader("rate limited")),
			}, nil
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("ok")),
		}, nil
	}, WithMaxAttempts(3), WithInitialDelay(time.Millisecond))

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
	resp.Body.Close()
}

func TestDoWithResponse_Retries5xx(t *testing.T) {
	calls := 0
	resp, err := DoWithResponse(func() (*http.Response, error) {
		calls++
		if calls < 3 {
			return &http.Response{
				StatusCode: 502,
				Body:       io.NopCloser(strings.NewReader("bad gateway")),
			}, nil
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("ok")),
		}, nil
	}, WithMaxAttempts(3), WithInitialDelay(time.Millisecond))

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	resp.Body.Close()
}

func TestDoWithResponse_NoRetryOn4xx(t *testing.T) {
	calls := 0
	resp, err := DoWithResponse(func() (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(strings.NewReader("bad request")),
		}, nil
	}, WithMaxAttempts(3), WithInitialDelay(time.Millisecond))

	if err != nil {
		t.Fatalf("expected nil error for 400 (not retryable at HTTP level), got %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call (no retry for 400), got %d", calls)
	}
	resp.Body.Close()
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"permanent", fmt.Errorf("bad input"), false},
		{"HTTP 400", &HTTPError{StatusCode: 400}, false},
		{"HTTP 404", &HTTPError{StatusCode: 404}, false},
		{"HTTP 429", &HTTPError{StatusCode: 429}, true},
		{"HTTP 500", &HTTPError{StatusCode: 500}, true},
		{"HTTP 502", &HTTPError{StatusCode: 502}, true},
		{"HTTP 503", &HTTPError{StatusCode: 503}, true},
		{"net.OpError", &net.OpError{Op: "read", Err: fmt.Errorf("timeout")}, true},
		{"ECONNRESET", syscall.ECONNRESET, true},
		{"ECONNREFUSED", syscall.ECONNREFUSED, true},
		{"EPIPE", syscall.EPIPE, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRetryable(tt.err)
			if got != tt.want {
				t.Errorf("IsRetryable(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
