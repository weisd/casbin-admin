package jwt

import (
	"net"
	"net/http"
	"strings"
	"time"
)

type jwtExtractor func(r *http.Request) string

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(r *http.Request) string {
		auth := r.Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:]
		}
		return ""
	}
}

// jwtFromQuery returns a `jwtExtractor` that extracts token from the query string.
func jwtFromQuery(param string) jwtExtractor {
	return func(r *http.Request) string {
		token := r.URL.Query().Get(param)
		if token == "" {
			return ""
		}
		return token
	}
}

// jwtFromCookie returns a `jwtExtractor` that extracts token from the named cookie.
func jwtFromCookie(name string) jwtExtractor {
	return func(r *http.Request) string {
		cookie, err := r.Cookie(name)
		if err != nil {
			return ""
		}

		return cookie.Value
	}
}

// RealIP RealIP
func RealIP(r *http.Request) string {
	ra := r.RemoteAddr
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		ra = strings.Split(ip, ", ")[0]
	} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}

// MinNight MinNight
func MinNight(t ...time.Time) time.Time {
	n := time.Now()
	if len(t) > 0 {
		n = t[0]
	}
	return time.Date(n.Year(), n.Month(), n.Day(), 3, 0, 0, 0, time.Local)
}
