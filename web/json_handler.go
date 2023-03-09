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
