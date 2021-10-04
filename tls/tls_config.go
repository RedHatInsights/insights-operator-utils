package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// NewTLSConfig create a TLS configuration from a certificate path. This can be
// used with Sarama for example.
func NewTLSConfig(certPath string) (*tls.Config, error) {
	if certPath == "" {
		return nil, fmt.Errorf("no cert path provided. Skip")
	}
	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{},
		MinVersion:   tls.VersionTLS12,
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	if caCert == nil {
		return nil, fmt.Errorf("pointer to new CertPool is nil")
	}

	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("error appending the specified certificate")
	}
	tlsConfig.RootCAs = caCertPool
	return &tlsConfig, err
}
