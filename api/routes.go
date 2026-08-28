package api

import (
	"net/http"

	"coursechain/workflow"
)

func NewRouter(service *workflow.Service) http.Handler {
	h := NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.Health)
	mux.HandleFunc("/api/v1/summary", h.Summary)
	mux.HandleFunc("/api/v1/timeline", h.Timeline)
	mux.HandleFunc("/api/v1/submissions", h.Submissions)
	mux.HandleFunc("/api/v1/submissions/", h.SubmissionByID)
	mux.HandleFunc("/api/v1/users", h.Users)
	return withJSON(withRecovery(mux))
}

func withJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
