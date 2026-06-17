package middleware

import (
	"log"
	"net/http"

	"sklbz-ng/config"
)

// AuthMiddleware wraps the mTLS middleware for write operations
func AuthMiddleware(mtlsConfig *config.MTLSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for GET, HEAD, OPTIONS requests
			if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			// For POST, PUT, DELETE, PATCH - require mTLS authentication
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" || r.Method == "PATCH" {
				// If mTLS is not enabled, we can still check for other auth methods
				if !mtlsConfig.Enabled {
					// For now, if mTLS is disabled, we allow the request
					// In production, you might want to require another form of auth
					log.Printf("Warning: mTLS is disabled, allowing %s request without authentication", r.Method)
					next.ServeHTTP(w, r)
					return
				}

				// Check for client certificate
				if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
					w.Header().Set("WWW-Authenticate", "mTLS")
					http.Error(w, "Client certificate required for write operations", http.StatusUnauthorized)
					return
				}

				// Verify the certificate is valid (already done by MTLSMiddleware)
				// But we can add additional checks here if needed
				clientCN := config.GetClientCommonName(r)
				if clientCN == "" {
					http.Error(w, "Invalid client certificate", http.StatusUnauthorized)
					return
				}

				// Log the authenticated request
				log.Printf("mTLS Auth: User '%s' making %s request to %s", 
					clientCN, r.Method, r.URL.Path)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireMTLS creates a middleware that requires mTLS for all requests
func RequireMTLS(mtlsConfig *config.MTLSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !mtlsConfig.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
				w.Header().Set("WWW-Authenticate", "mTLS")
				http.Error(w, "Client certificate required", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// WriteOnlyAuth creates a middleware that only requires auth for write operations
func WriteOnlyAuth(mtlsConfig *config.MTLSConfig) func(http.Handler) http.Handler {
	return AuthMiddleware(mtlsConfig)
}
