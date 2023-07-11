package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/evertras/rcc/pkg/badge"
)

func v0HandlerBadgeCoverage(getter coverageValueGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get(queryParamKey)

		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Missing key"))
			return
		}

		getThreshold := func(queryParam string, defaultValue int) (int, error) {
			raw := r.URL.Query().Get(queryParam)

			if raw == "" {
				return defaultValue, nil
			}

			parsed, err := strconv.Atoi(strings.TrimRight(raw, "%"))

			if err != nil {
				return 0, fmt.Errorf("could not parse value")
			}

			if parsed < 0 || parsed > 100 {
				return 0, fmt.Errorf("threshold value must be between 0-100")
			}

			return parsed, nil
		}

		thresholdOrange100, err := getThreshold(queryParamThresholdOrange100, 80)

		if err != nil {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(fmt.Sprintf("Bad orange threshold value: %s", err.Error())))
			return
		}

		thresholdRed100, err := getThreshold(queryParamThresholdRed100, 50)

		if err != nil {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(fmt.Sprintf("Bad red threshold value: %s", err.Error())))
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

		if value1000 < thresholdRed100*10 {
			color = badge.ColorRed
		} else if value1000 < thresholdOrange100*10 {
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
