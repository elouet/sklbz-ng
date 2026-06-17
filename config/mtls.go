package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// MTLSConfig holds the mTLS configuration
type MTLSConfig struct {
	Enabled          bool
	ClientAuthType   tls.ClientAuthType
	CertPool         *x509.CertPool
	MinVersion       uint16
	CipherSuites     []uint16
}

// DefaultMTLSConfig returns a default mTLS configuration
func DefaultMTLSConfig() *MTLSConfig {
	return &MTLSConfig{
		Enabled:        false,
		ClientAuthType: tls.RequireAndVerifyClientCert,
		CertPool:       x509.NewCertPool(),
		MinVersion:     tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

// LoadMTLSConfig loads mTLS configuration from environment or files
func LoadMTLSConfig(certFile, keyFile, caCertFile string) (*MTLSConfig, error) {
	config := DefaultMTLSConfig()

	// Check if mTLS is enabled via environment variable
	if os.Getenv("MTLS_ENABLED") == "true" {
		config.Enabled = true
	} else {
		// If cert files are provided, enable mTLS
		if certFile != "" && keyFile != "" && caCertFile != "" {
			config.Enabled = true
		}
	}

	if !config.Enabled {
		return config, nil
	}

	// Load CA certificate for client verification
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}

	if !config.CertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	log.Printf("mTLS enabled: client certificate verification required")
	log.Printf("CA certificate loaded from: %s", caCertFile)

	return config, nil
}

// MTLSMiddleware creates a middleware that checks for client certificates
func MTLSMiddleware(config *MTLSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If mTLS is not enabled, just pass through
			if !config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Check if the request has a valid client certificate
			if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
				http.Error(w, "Client certificate required", http.StatusUnauthorized)
				return
			}

			// Verify the client certificate is in our CA pool
			cert := r.TLS.PeerCertificates[0]
			if _, err := cert.Verify(x509.VerifyOptions{
				Roots:         config.CertPool,
				Intermediates: x509.NewCertPool(),
			}); err != nil {
				http.Error(w, "Invalid client certificate: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Check if certificate is expired
			// Parse the Date header, default to current time if not present
			dateHeader := r.Header.Get("Date")
			var requestTime time.Time
			if dateHeader != "" {
				var err error
				requestTime, err = http.ParseTime(dateHeader)
				if err != nil {
					log.Printf("Error parsing Date header: %v", err)
					requestTime = time.Now()
				}
			} else {
				requestTime = time.Now()
			}
			
			if cert.NotAfter.Before(requestTime) {
				http.Error(w, "Client certificate expired", http.StatusUnauthorized)
				return
			}

			// Add certificate info to request context for logging/auditing
			r.Header.Set("X-Client-CN", cert.Subject.CommonName)
			r.Header.Set("X-Client-Cert-Issuer", cert.Issuer.CommonName)

			log.Printf("mTLS: Client authenticated - CN: %s, Issuer: %s", 
				cert.Subject.CommonName, cert.Issuer.CommonName)

			next.ServeHTTP(w, r)
		})
	}
}

// CreateMTLSListener creates a TLS listener with mTLS support
func CreateMTLSListener(addr, certFile, keyFile string, mtlsConfig *MTLSConfig) (*http.Server, error) {
	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate/key: %v", err)
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   mtlsConfig.ClientAuthType,
		ClientCAs:    mtlsConfig.CertPool,
		MinVersion:   mtlsConfig.MinVersion,
		CipherSuites: mtlsConfig.CipherSuites,
	}

	// Create server with TLS config
	server := &http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
	}

	return server, nil
}

// GetClientCommonName extracts the client's common name from the certificate
func GetClientCommonName(r *http.Request) string {
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		return ""
	}
	return r.TLS.PeerCertificates[0].Subject.CommonName
}

// IsMTLSRequest checks if the request was made with a valid client certificate
func IsMTLSRequest(r *http.Request) bool {
	return r.TLS != nil && len(r.TLS.PeerCertificates) > 0
}
