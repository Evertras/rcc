package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/evertras/rcc/pkg/badge"
)

func v0HandlerBadgeCoverage(getter coverageValueGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")

		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing key"))
			return
		}

		value1000, err := getter.GetValue1000(r.Context(), key)

		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				// Return 200, let the badge visually show that we didn't find it
				w.Write([]byte(badge.GenerateCoverageUnknownSVG()))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("ERROR: Failed to get data for key %q: %s", key, err.Error())
			return
		}

		b, err := badge.GenerateCoverageSVG(value1000)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("ERROR: Failed to generate badge from value %d: %s", value1000, err.Error())
			return
		}

		w.Write([]byte(b))
	}
}
