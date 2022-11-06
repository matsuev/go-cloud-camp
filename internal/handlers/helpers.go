package handlers

import (
	"go-cloud-camp/internal/common"
	"net/http"
	"strconv"
)

// setContentTypeJSON function
func (h *AppHandlers) setContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// getServiceAndVersion function
func (h *AppHandlers) getServiceAndVersion(r *http.Request) (string, int, error) {
	requestQuery := r.URL.Query()

	service := requestQuery.Get("service")
	if service == common.EMPTY_STRING {
		return common.EMPTY_STRING, 0, common.ErrEmptyServiceName
	}

	// Если параметр version не задан, или это не число,
	// тогда version = 0, это значит выбрать последнюю версию конфига
	version, _ := strconv.Atoi(requestQuery.Get("version"))

	return service, version, nil
}
