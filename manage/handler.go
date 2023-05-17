package manage

import (
	"circa/config"
	"encoding/json"
	"net/http"
)

type configManage struct {
	c *config.Config
}

func newConfigManage(c *config.Config) *configManage {
	return &configManage{c: c}
}

func (cm *configManage) Sync(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if r.Method == http.MethodPost {
		if err := cm.c.Sync(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(&errorResponse{"ERROR_SYNC_CONFIG", "Can't sync config"})
		}
		return
	}
	writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, err *errorResponse) {
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(err)
}
