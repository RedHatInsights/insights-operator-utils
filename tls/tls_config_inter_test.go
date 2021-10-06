package tlsutil

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
