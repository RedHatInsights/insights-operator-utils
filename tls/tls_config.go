package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type tlsConfigGetter interface {
	newCertPoolGetter() *x509.CertPool
	appendCertsFromPEMGetter(caCertPool *x509.CertPool, pemCerts []byte) (ok bool)
}

type tlsConfigX509Getter struct{}

// NewTLSConfig create a TLS configuration from a certificate path. This can be
// used with Sarama for example.
func NewTLSConfig(certPath string) (*tls.Config, error) {
	tcg := new(tlsConfigX509Getter)
	return newTLSConfig(tcg, certPath)
}

func (t tlsConfigX509Getter) newCertPoolGetter() *x509.CertPool {
	return x509.NewCertPool()
}

func (t tlsConfigX509Getter) appendCertsFromPEMGetter(caCertPool *x509.CertPool, pemCerts []byte) (ok bool) {
	return caCertPool.AppendCertsFromPEM(pemCerts)
}

func newTLSConfig(t tlsConfigGetter, certPath string) (*tls.Config, error) {
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
	caCertPool := t.newCertPoolGetter()
	if caCertPool == nil {
		return nil, fmt.Errorf("pointer to new CertPool is nil")
	}

	ok := t.appendCertsFromPEMGetter(caCertPool, caCert)
	if !ok {
		return nil, fmt.Errorf("error appending the specified certificate")
	}
	tlsConfig.RootCAs = caCertPool
	return &tlsConfig, err
}
