package health

import "net/http"

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if err := h.postgres.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("fail"))

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))

	return
}
