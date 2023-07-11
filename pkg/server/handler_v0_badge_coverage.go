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
			_, _ = w.Write([]byte("Missing key"))
			return
		}

		value1000, err := getter.GetValue1000(r.Context(), key)

		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				// Return 200, let the badge visually show that we didn't find it
				w.Header().Add("Content-Type", "image/svg+xml")
				_, _ = w.Write([]byte(badge.GenerateCoverageUnknownSVG()))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("ERROR: Failed to get data for key %q: %s", key, err.Error())
			return
		}

		color := badge.ColorGreen

		if value1000 < 500 {
			color = badge.ColorRed
		} else if value1000 < 800 {
			color = badge.ColorOrange
		}

		b, err := badge.GenerateCoverageSVG("coverage", value1000, color)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("ERROR: Failed to generate badge from value %d: %s", value1000, err.Error())
			return
		}

		w.Header().Add("Content-Type", "image/svg+xml")
		_, _ = w.Write([]byte(b))
	}
}
