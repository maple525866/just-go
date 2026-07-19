package profiler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterAddsPprofIndex(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/debug/pprof/", http.NoBody))
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAllocateHotHeapReturnsAllocatedBytes(t *testing.T) {
	got := AllocateHotHeap(3, 1024)
	if got != 3072 {
		t.Fatalf("allocated bytes = %d", got)
	}
}
