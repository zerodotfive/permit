package activate

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zerodotfive/permit/internal/permits"
)

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.Header.Get("X-Forwarded-For")

	w.Header().Set("Content-Type", "text/plain")

	provider := permits.NewProvider(r.Context(), h.postgres)

	rawToken := r.PathValue("token")
	if len(rawToken) != 43 {
		h.logger.Printf("host=%s remoteIP=%s activate.post: wrong token length %d", r.Host, remoteIP, len(rawToken))
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	token := make([]byte, 32)
	_, err := base64.RawURLEncoding.Decode(token, []byte(rawToken))
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s activate.post: token error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	permit, err := provider.ActivateToken(token)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s activate.post: token error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	sessionToken := base64.RawURLEncoding.EncodeToString(permit.SessionToken)
	chainToken := base64.RawURLEncoding.EncodeToString(permit.ChainToken)

	targetSite := strings.Trim(permit.TargetSite, "/")
	locationChain := strings.Trim(h.locationChain, "/")

	cts := targetSite + "/" + locationChain + "/" + chainToken
	ctu, err := url.Parse(cts)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s activate.post: chain site address wrong: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	esu, err := url.Parse(permit.EmbedSite)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s activate.post: target site address wrong: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	if esu.Host != r.Host {
		h.logger.Printf("host=%s remoteIP=%s activate.post: target site and host header mismatch: %s, %s", r.Host, remoteIP, esu.Host, r.Host)
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "permit",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(*permit.SessionTokenUseNotAfter).Seconds()),
	})

	h.logger.Printf("host=%s remoteIP=%s activate.post: success", r.Host, remoteIP)
	http.Redirect(w, r, ctu.String(), http.StatusFound)
	return
}
