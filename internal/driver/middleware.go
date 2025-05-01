package driver

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// BasicAuth middleware checks the username and password from the Authorization header
func BasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Basic ") {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(auth[len("Basic "):])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) != 2 || parts[0] != username || parts[1] != password {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
