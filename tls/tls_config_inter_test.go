// Copyright 2021 Red Hat, Inc
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

package tlsutil

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/tls/tls_config_inter_test.html

import (
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewTLSConfig tests the NewTLSConfig method
func TestNewTLSConfigWithMock(t *testing.T) {
	var certPath = "../testdata/cert.pem"
	t.Run("bad x509.CertPool", func(t *testing.T) {
		_, err := newTLSConfig(tlsConfigX509GetterWithNilInNewCertPoolGetter{}, certPath)
		assert.Error(t, err)
	})
	t.Run("bad appendCertsFromPEMGetter", func(t *testing.T) {
		_, err := newTLSConfig(tlsConfigX509GetterWithFalseInAppendCertsFromPEMGetter{}, certPath)
		assert.Error(t, err)
	})
}

type tlsConfigX509GetterWithNilInNewCertPoolGetter struct{}

func (t tlsConfigX509GetterWithNilInNewCertPoolGetter) newCertPoolGetter() *x509.CertPool {
	return nil
}

func (t tlsConfigX509GetterWithNilInNewCertPoolGetter) appendCertsFromPEMGetter(caCertPool *x509.CertPool, pemCerts []byte) (ok bool) {
	return caCertPool.AppendCertsFromPEM(pemCerts)
}

type tlsConfigX509GetterWithFalseInAppendCertsFromPEMGetter struct{}

func (t tlsConfigX509GetterWithFalseInAppendCertsFromPEMGetter) newCertPoolGetter() *x509.CertPool {
	return x509.NewCertPool()
}

func (t tlsConfigX509GetterWithFalseInAppendCertsFromPEMGetter) appendCertsFromPEMGetter(caCertPool *x509.CertPool, pemCerts []byte) (ok bool) {
	return false
}
