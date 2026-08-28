package api

import (
	"net/http"
)

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodAllowed(w, "GET")
		return
	}
	report, err := h.service.Report(r.Context(), r.URL.Query().Get("course"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) Timeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodAllowed(w, "GET")
		return
	}
	id := r.URL.Query().Get("record_id")
	items, err := h.service.Timeline(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"record_id": id, "timeline": items})
}
