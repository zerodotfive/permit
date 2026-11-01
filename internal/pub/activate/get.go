package activate

import (
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/zerodotfive/permit/internal/permits"
	"github.com/zerodotfive/permit/ui"
)

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.Header.Get("X-Forwarded-For")

	w.Header().Set("Content-Type", "text/html")

	provider := permits.NewProvider(r.Context(), h.postgres)

	rawToken := r.PathValue("token")
	if len(rawToken) != 43 {
		h.logger.Printf("host=%s remoteIP=%s activate.get: wrong token length %d", r.Host, remoteIP, len(rawToken))
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	token := make([]byte, 32)
	_, err := base64.RawURLEncoding.Decode(token, []byte(rawToken))
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s activate.get: token error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	permit, err := provider.InquireForToken(token)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s activate.get: host=%s remoteIP=%s: inquire token error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	esu, err := url.Parse(permit.EmbedSite)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s activate.get: target site address wrong: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	if esu.Host != r.Host {
		h.logger.Printf("host=%s remoteIP=%s activate.get: target site and host header mismatch: %s, %s", r.Host, remoteIP, esu.Host, r.Host)
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	h.logger.Printf("host=%s remoteIP=%s activate.get: success", r.Host, remoteIP)
	fileBytes, _ := ui.GetEmbedHTMLBytes("client.html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(fileBytes)
	return
}
