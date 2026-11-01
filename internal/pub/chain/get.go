package chain

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zerodotfive/permit/internal/permits"
)

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.Header.Get("X-Forwarded-For")

	w.Header().Set("Content-Type", "text/plain")

	provider := permits.NewProvider(r.Context(), h.postgres)

	rawToken := r.PathValue("token")
	if len(rawToken) != 43 {
		h.logger.Printf("host=%s remoteIP=% chain.get: wrong token length %d", r.Host, remoteIP, len(rawToken))
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	token := make([]byte, 32)
	_, err := base64.RawURLEncoding.Decode(token, []byte(rawToken))
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s chain.get: token error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	permit, err := provider.ChainToken(token)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s chain.get: token error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	sessionToken := base64.RawURLEncoding.EncodeToString(permit.SessionToken)

	targetSite := strings.Trim(permit.TargetSite, "/")
	locationTarget := strings.TrimLeft(permit.TargetLocation, "/")

	tts := targetSite + "/" + locationTarget
	ttu, err := url.Parse(tts)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s chain.get: token error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	csu, err := url.Parse(permit.TargetSite)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s chain.get: target site address wrong: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusNotFound)
		return
	}

	if csu.Host != r.Host {
		h.logger.Printf("host=%s remoteIP=%s chain.get: target site and host header mismatch: %s, %s", r.Host, remoteIP, csu.Host, r.Host)
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

	h.logger.Printf("host=%s remoteIP=%s chain.get: success", r.Host, remoteIP)
	http.Redirect(w, r, ttu.String(), http.StatusFound)
	return
}
