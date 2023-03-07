package web

import (
	"encoding/json"
	"net/http"
)

type JsonHandler struct{}

func (h *JsonHandler) readJSON(r *http.Request, val any) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return err
	}
	return nil
}

// func (h *JsonHandler) writeJSON(w http.ResponseWriter, val any, code int) error {
// 	w.Header().Add("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(val); err != nil {
// 		return err
// 	}
// 	w.WriteHeader(code)
// 	return nil
// }
