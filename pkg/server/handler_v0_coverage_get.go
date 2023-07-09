package server

import (
	"context"
	"net/http"
)

type coverageValueGetter interface {
	GetValue1000(ctx context.Context, key string) (int, error)
}

func v0HandlerCoverageGet(getter coverageValueGetter) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			queryVals := r.URL.Query()

			key := queryVals.Get("key")
			if key == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Missing key"))
				return
			}
		},
	)
}
