package server

import "net/http"

func handlerHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For now, we're always healthy, so just let it write 200 OK naturally
	}
}
