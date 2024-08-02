package profiler

import (
	"net/http"
	npprof "net/http/pprof"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", npprof.Index)
	r.Get("/cmdline", npprof.Cmdline)
	r.Get("/profile", npprof.Profile)
	r.Get("/trace", npprof.Trace)
	r.Get("/allocs", npprof.Handler("allocs").ServeHTTP)
	r.Get("/block", npprof.Handler("block").ServeHTTP)
	r.Get("/goroutine", npprof.Handler("goroutine").ServeHTTP)
	r.Get("/heap", npprof.Handler("heap").ServeHTTP)
	r.Get("/mutex", npprof.Handler("mutex").ServeHTTP)
	r.Get("/threadcreate", npprof.Handler("threadcreate").ServeHTTP)

	return r
}
