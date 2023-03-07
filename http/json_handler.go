package http

import "net/http"

type JsonHandler struct{}

func (h *JsonHandler) readJSON(r *http.Request, val any) error {
	return nil
}

func (h *JsonHandler) writeJSON(w http.ResponseWriter, val any) error {
	return nil
}
