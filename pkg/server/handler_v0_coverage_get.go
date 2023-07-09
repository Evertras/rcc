package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
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

			val1000, err := getter.GetValue1000(r.Context(), key)

			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("ERROR: Failed to get coverage for key %q: %s", key, err.Error())
				return
			}

			fmt.Fprintf(w, "%d.%d", val1000/10, val1000%10)
		},
	)
}
