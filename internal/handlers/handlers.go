package handlers

import (
	"encoding/json"
	"errors"
	"go-cloud-camp/internal/common"
	"go-cloud-camp/internal/logging"
	"go-cloud-camp/internal/storage"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	configURL = "/config"
)

// AppHandlers struct
type AppHandlers struct {
	Log     *logging.Logger
	Storage *storage.AppStorage
}

// Create function
func Create(l *logging.Logger, s *storage.AppStorage) *AppHandlers {
	return &AppHandlers{
		Log:     l,
		Storage: s,
	}
}

// Register function
func (h *AppHandlers) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, configURL, h.Get)
	router.HandlerFunc(http.MethodPost, configURL, h.Post)
	router.HandlerFunc(http.MethodPut, configURL, h.Put)
	router.HandlerFunc(http.MethodDelete, configURL, h.Delete)
}

// Get function
func (h *AppHandlers) Get(w http.ResponseWriter, r *http.Request) {
	service, version, err := h.getServiceAndVersion(r)
	if err != nil {
		// Error 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.Storage.Read(service, version)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			// Error 404
			w.WriteHeader(http.StatusNotFound)
		default:
			// Error 500
			w.WriteHeader(http.StatusInternalServerError)
		}

		h.LogInfoRequestDetails("GET request aborted with error", err, r)

		return
	}

	h.setContentTypeJSON(w)
	if _, err = w.Write(result); err != nil {
		h.LogDebugRequestDetails("http.ResponseWriter was called with an error", err, r)
	}

	h.LogRequest("GET request completed", r)
}

// Post function
func (h *AppHandlers) Post(w http.ResponseWriter, r *http.Request) {
	postData := &common.RequestData{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(postData); err != nil {
		h.LogInfoRequestDetails("POST request aborted with error", err, r)
		// Error 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Storage.Create(postData); err != nil {
		switch {
		case errors.Is(err, common.ErrNotValidJsonData):
			// Error 400
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, common.ErrAlreadyCreated):
			// Error 403
			w.WriteHeader(http.StatusForbidden)
		default:
			// Error 500
			w.WriteHeader(http.StatusInternalServerError)
		}

		h.LogInfoRequestDetails("POST request aborted with error", err, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.LogRequest("POST request completed", r)
}

// Put function
func (h *AppHandlers) Put(w http.ResponseWriter, r *http.Request) {
	postData := &common.RequestData{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(postData); err != nil {
		h.LogInfoRequestDetails("GET request aborted with error", err, r)
		// Error 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Storage.Update(postData); err != nil {
		switch {
		case errors.Is(err, common.ErrNotValidJsonData):
			// Error 400
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, common.ErrServiceNotFound):
			// Error 404
			w.WriteHeader(http.StatusNotFound)
		default:
			// Error 500
			w.WriteHeader(http.StatusInternalServerError)
		}

		h.LogInfoRequestDetails("PUT request aborted with error", err, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	h.LogRequest("PUT request completed", r)
}

// Delete function
func (h *AppHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	service, version, err := h.getServiceAndVersion(r)
	if err != nil {
		// Error 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Storage.Delete(service, version); err != nil {
		switch {
		case errors.Is(err, common.ErrConfigIsUsed):
			// Error 403
			w.WriteHeader(http.StatusForbidden)
		case errors.Is(err, common.ErrServiceNotFound):
			// Error 404
			w.WriteHeader(http.StatusNotFound)
		default:
			// Error 500
			w.WriteHeader(http.StatusInternalServerError)
		}

		h.LogInfoRequestDetails("DELETE request aborted with error", err, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	h.LogRequest("DELETE request completed", r)
}

// LogRequest function
func (h *AppHandlers) LogRequest(msg string, r *http.Request) {
	h.Log.Infow(msg,
		"remote_addr", r.RemoteAddr,
		"request_uri", r.RequestURI,
	)
}

// LogInfoRequestDetails function
func (h *AppHandlers) LogInfoRequestDetails(msg string, err error, r *http.Request) {
	h.Log.Infow(msg,
		"error", err,
		"reamote_addr", r.RemoteAddr,
		"request_uri", r.RequestURI,
	)
}

// LogDebugwRequestDetails function
func (h *AppHandlers) LogDebugRequestDetails(msg string, err error, r *http.Request) {
	h.Log.Debugw(msg,
		"error", err,
		"reamote_addr", r.RemoteAddr,
		"request_uri", r.RequestURI,
	)
}
