package auth

import (
	"net/http"
	"time"
)

// isSecureRequest reports whether the client reached us over HTTPS, directly
// or through a TLS-terminating proxy such as dokku's nginx.
func isSecureRequest(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

type authCookieValues struct {
	name  string
	value string
	// domain   string
	httpOnly bool
	lifetime time.Duration
	secure   bool
}

func makeAuthCookie(values authCookieValues) *http.Cookie {
	expiresAt := time.Now().Add(values.lifetime)

	cookie := &http.Cookie{
		Name:     values.name,
		Value:    values.value,
		Expires:  expiresAt,
		HttpOnly: values.httpOnly,
		Path:     "/",
		// Domain:   values.domain,
		MaxAge:   int(values.lifetime.Seconds()),
		Secure:   values.secure,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}

func clearAuthCookie(name string, httpOnly bool, domain string, secure bool) *http.Cookie {
	vals := authCookieValues{
		name:  name,
		value: "",
		// domain:   domain,
		httpOnly: httpOnly,
		lifetime: time.Second,
		secure:   secure,
	}
	cookie := makeAuthCookie(vals)
	return cookie
}
