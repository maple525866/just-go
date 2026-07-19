package profiler

import (
	"net/http"
	"net/http/pprof"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func AllocateHotHeap(blocks int, blockSize int) int {
	if blocks <= 0 || blockSize <= 0 {
		return 0
	}
	data := make([][]byte, blocks)
	total := 0
	for i := range data {
		data[i] = make([]byte, blockSize)
		data[i][0] = byte(i)
		total += len(data[i])
	}
	return total
}
