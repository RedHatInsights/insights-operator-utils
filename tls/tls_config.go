// Copyright 2021, 2022 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package tlsutil contains helper function to create TLS configurations
package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
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
	caCert, err := os.ReadFile(filepath.Clean(certPath))
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
