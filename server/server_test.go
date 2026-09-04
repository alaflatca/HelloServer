package server

import (
	"helloServer/cache"
	"helloServer/event"
	"helloServer/measure"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMetricsCollectionRoute(t *testing.T) {
	cache.Set("l", &measure.Measure{})
	svr := New()
	svr.registerRoutes()

	resp := testRequest(t, svr, http.MethodGet, "/api/metrics", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestMetricsAllRouteIsNoLongerUsed(t *testing.T) {
	cache.Set("l", &measure.Measure{})
	svr := New()
	svr.registerRoutes()

	resp := testRequest(t, svr, http.MethodGet, "/api/metrics/all", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestAgentPeriodRoute(t *testing.T) {
	periodCh := make(chan time.Duration, 1)
	event.Subscribe("period", func(data interface{}) {
		period, ok := data.(time.Duration)
		if ok {
			periodCh <- period
		}
	})

	svr := New()
	svr.registerRoutes()

	resp := testRequest(t, svr, http.MethodPut, "/api/agent/period", `{"period":3}`)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	select {
	case period := <-periodCh:
		if period != 3*time.Second {
			t.Fatalf("period = %s, want 3s", period)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("period event was not published")
	}
}

func TestOldPeriodRouteIsNotRegistered(t *testing.T) {
	svr := New()
	svr.registerRoutes()

	resp := testRequest(t, svr, http.MethodPut, "/api/period", `{"period":3}`)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func testRequest(t *testing.T, svr *server, method, path, body string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := svr.app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
