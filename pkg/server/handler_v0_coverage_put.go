package server

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type CoverageValueStorer interface {
	// StoreValue1000 stores a value as an integer that includes 1 decimal
	StoreValue1000(ctx context.Context, key string, value1000 int) error
}

func v0HandlerCoveragePut(storer CoverageValueStorer) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			queryVals := r.URL.Query()

			key := queryVals.Get("key")
			if key == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Missing key"))
				return
			}

			value100Raw := queryVals.Get("value100")
			if value100Raw == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Missing value100"))
				return
			}

			value100Raw = strings.TrimSuffix(value100Raw, "%")

			parsedValue1000, err := parseValue1000(value100Raw)

			if err != nil {
				log.Printf("Failed to parse %q: %s", value100Raw, err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			err = storer.StoreValue1000(r.Context(), key, parsedValue1000)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Something went wrong when storing the value"))
				return
			}
		},
	)
}
