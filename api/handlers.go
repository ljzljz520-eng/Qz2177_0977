package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"coursechain/domain"
	"coursechain/query"
	"coursechain/workflow"
)

type Handler struct {
	service *workflow.Service
	query   *query.Engine
}

func NewHandler(service *workflow.Service) *Handler {
	return &Handler{service: service, query: query.New(service.Store())}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "course": h.service.Course()})
}

func (h *Handler) Submissions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var input domain.Submission
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		actor := r.Header.Get("X-Actor-ID")
		if actor == "" {
			actor = input.StudentID
		}
		record, err := h.service.Submit(r.Context(), input, actor)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, record)
	case http.MethodGet:
		filter, err := parseFilter(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		page, err := h.query.Search(r.Context(), filter)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, page)
	default:
		methodAllowed(w, "GET, POST")
	}
}

func (h *Handler) SubmissionByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/submissions/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		record, err := h.service.Store().GetRecord(r.Context(), id)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, record)
	case http.MethodPost:
		if r.URL.Query().Get("action") == "notify" {
			actor := r.Header.Get("X-Actor-ID")
			if err := h.service.Notify(r.Context(), id, actor, r.URL.Query().Get("detail")); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			writeJSON(w, http.StatusAccepted, map[string]string{"status": "notified"})
			return
		}
		writeError(w, http.StatusBadRequest, "unsupported action")
	default:
		methodAllowed(w, "GET, POST")
	}
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodAllowed(w, "POST")
		return
	}
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.service.RegisterUser(r.Context(), user); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func parseFilter(r *http.Request) (domain.QueryFilter, error) {
	values := r.URL.Query()
	limit, err := parseNumber(values.Get("limit"), 50)
	if err != nil {
		return domain.QueryFilter{}, err
	}
	offset, err := parseNumber(values.Get("offset"), 0)
	if err != nil {
		return domain.QueryFilter{}, err
	}
	filter := domain.QueryFilter{Course: values.Get("course"), StudentID: values.Get("student_id"), Search: values.Get("search"), Tag: values.Get("tag"), Limit: limit, Offset: offset}
	if raw := values.Get("status"); raw != "" {
		filter.Status = domain.Status(raw)
		if err := domain.ValidateStatus(filter.Status); err != nil {
			return domain.QueryFilter{}, err
		}
	}
	return filter, nil
}

func parseNumber(raw string, fallback int) (int, error) {
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, context.Canceled
	}
	return value, nil
}
