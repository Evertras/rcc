package server

import "net/http"

func v0HandlerCoveragePut() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")

			if key == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Missing key"))
				return
			}
		},
	)
}
