package verify

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zerodotfive/permit/internal/permits"
)

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.Header.Get("X-Forwarded-For")

	w.Header().Set("Content-Type", "text/plain")
	provider := permits.NewProvider(r.Context(), h.postgres)

	var rawSessionToken []byte

	cookieSessionToken, err := r.Cookie("permit")
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s verify.get: session token get error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusForbidden)
		return
	}

	if cookieSessionToken != nil && len(cookieSessionToken.Value) > 0 {
		rawSessionToken = []byte(cookieSessionToken.Value)
	}

	if len(rawSessionToken) == 0 {
		h.logger.Printf("host=%s remoteIP=%s verify.get: no session token found", r.Host, remoteIP)
		http.Error(w, "token error", http.StatusForbidden)
		return
	}

	if len(rawSessionToken) != 43 {
		h.logger.Printf("host=%s remoteIP=%s verify.get: wrong token length %d", r.Host, remoteIP, len(rawSessionToken))
		http.Error(w, "token error", http.StatusForbidden)
		return
	}

	sessionToken := make([]byte, 32)
	_, err = base64.RawURLEncoding.Decode(sessionToken, rawSessionToken)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s verify.get: token error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusForbidden)
		return
	}

	permit, err := provider.VerifySession(sessionToken)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s verify.get: permit use error: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "permit use error", http.StatusForbidden)
		return
	}

	esu, err := url.Parse(permit.EmbedSite)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s verify.get: embed site address wrong: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusForbidden)
		return
	}

	tsu, err := url.Parse(permit.TargetSite)
	if err != nil {
		h.logger.Printf("host=%s remoteIP=%s verify.get: target site address wrong: %s", r.Host, remoteIP, err.Error())
		http.Error(w, "token error", http.StatusForbidden)
		return
	}

	if esu.Host != r.Host && tsu.Host != r.Host {
		h.logger.Printf("host=%s remoteIP=%s verify.get: embed or target site mismatch: %s, %s", r.Host, remoteIP, esu.Host, tsu.Host)
		http.Error(w, "token error", http.StatusForbidden)
		return
	}

	h.logger.Printf("host=%s remoteIP=%s verify.get: success", r.Host, remoteIP)
	_, _ = fmt.Fprintf(w, "ok")
	return
}
