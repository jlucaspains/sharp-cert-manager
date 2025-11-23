package handlers

import (
	"net/http"
)

func (h Handlers) CORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", h.CORSOrigins)
}
