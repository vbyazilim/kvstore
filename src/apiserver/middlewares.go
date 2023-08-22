package apiserver

import (
	"log/slog"
	"net/http"
	"strings"
)

func httpLoggingMiddleware(l *slog.Logger, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		uri := r.URL.String()
		method := r.Method

		l.Info("http request", "method", method, "uri", uri)
	}

	return http.HandlerFunc(fn)
}

func appendSlashMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && !strings.HasSuffix(r.URL.Path, "/") {
			redirectURL := r.URL.Path + "/"
			if r.URL.RawQuery != "" {
				redirectURL += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
