package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *HTTPDServer) handleGetIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.Json(w, map[string]interface{}{"hello": "world"})
}
