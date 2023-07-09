package server

import (
	"log"
	"net/http"
	"strings"
)

func v0HandlerCoveragePut(storer coverageValueStorer) http.HandlerFunc {
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

			log.Printf("Setting coverage for %q to %d", key, parsedValue1000)

			err = storer.StoreValue1000(r.Context(), key, parsedValue1000)

			if err != nil {
				log.Println("Failed to store value:", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Something went wrong when storing the value"))
				return
			}

			log.Println("Successfully stored coverage value")

			w.WriteHeader(http.StatusOK)
		},
	)
}
