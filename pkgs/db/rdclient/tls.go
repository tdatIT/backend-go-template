package rdclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/tdatIT/backend-go/config"
)

// configureTLS configures TLS settings for Redis client based on config
func configureTLS(cfg *config.ServiceConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.Redis.TLS.InsecureSkipVerify, // #nosec G402 -- internal redis, private network, no public exposure
	}

	// Load CA certificate if provided
	if cfg.Redis.TLS.CertFilePath != "" {
		caCert, err := os.ReadFile(cfg.Redis.TLS.CertFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate from %s: %w", cfg.Redis.TLS.CertFilePath, err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate from %s", cfg.Redis.TLS.CertFilePath)
		}

		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}
