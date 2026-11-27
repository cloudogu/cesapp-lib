package ssl_test

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/cloudogu/cesapp-lib/ssl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestNewSSLGenerator(t *testing.T) {
	// when
	generator := ssl.NewSSLGenerator()

	// then
	require.NotNil(t, generator)
}

func Test_sslGenerator_GenerateSelfSignedCert(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		generator := ssl.NewSSLGenerator()

		// when
		cert, key, err := generator.GenerateSelfSignedCert("fqdn", "myces", 365, "de",
			"lower sachs", "brunswick", nil)

		// then
		require.NoError(t, err)
		assert.NotEmpty(t, cert)
		assert.NotEmpty(t, key)

		certs := splitPemCertificates(cert)
		assert.Equal(t, 2, len(certs))

		err = validateCert(certs[0])
		require.NoError(t, err)
		err = validateCert(certs[1])
		require.NoError(t, err)
		_, err = validatePEM(key)
		require.NoError(t, err)
	})

	t.Run("should not set fqdn as common name", func(t *testing.T) {
		// given
		generator := ssl.NewSSLGenerator()

		// when
		cert, _, err := generator.GenerateSelfSignedCert("my.ces.example.com", "myces", 365, "de",
			"lower sachs", "brunswick", nil)

		// then
		require.NoError(t, err)
		leaf := parseCert(t, splitPemCertificates(cert)[0])
		assert.Empty(t, leaf.Subject.CommonName)
		assert.Contains(t, leaf.DNSNames, "my.ces.example.com")
	})

	t.Run("should succeed with fqdn longer than 64 characters", func(t *testing.T) {
		// given
		generator := ssl.NewSSLGenerator()
		longFqdn := strings.Repeat("a", 60) + ".subdomain.example.com"
		require.Greater(t, len(longFqdn), 64)

		// when
		cert, _, err := generator.GenerateSelfSignedCert(longFqdn, "myces", 365, "de",
			"lower sachs", "brunswick", nil)

		// then
		require.NoError(t, err)
		leaf := parseCert(t, splitPemCertificates(cert)[0])
		assert.Empty(t, leaf.Subject.CommonName)
		assert.Contains(t, leaf.DNSNames, longFqdn)
		require.NoError(t, leaf.VerifyHostname(longFqdn))
	})
}

func parseCert(t *testing.T, certPEM string) *x509.Certificate {
	t.Helper()
	block, err := validatePEM(certPEM)
	require.NoError(t, err)
	parsed, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)
	return parsed
}

func splitPemCertificates(chain string) []string {
	sep := "-----BEGIN CERTIFICATE-----\n"
	var result []string
	split := strings.Split(chain, sep)
	for _, s := range split {
		if s == "" {
			continue
		}
		result = append(result, fmt.Sprintf("%s%s", sep, s))
	}
	return result
}

func validateCert(cert string) error {
	block, err := validatePEM(cert)
	_, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return assert.AnError
	}

	return nil
}

func validatePEM(pemStr string) (*pem.Block, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, assert.AnError
	}

	return block, nil
}
