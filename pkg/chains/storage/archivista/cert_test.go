package archivista

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"strings"
	"testing"
	"time"
)

// createCertificate is a helper that generates a PEM-encoded certificate, its parsed form, and a private key.
// If parent and parentKey are nil, the certificate is self-signed.
func createCertificate(t *testing.T, subject, issuer string, serial int64, notBefore time.Time, parent *x509.Certificate, parentKey *rsa.PrivateKey) (string, *x509.Certificate, *rsa.PrivateKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(serial),
		Subject: pkix.Name{
			CommonName: subject,
		},
		Issuer: pkix.Name{
			CommonName: issuer,
		},
		NotBefore: notBefore,
		NotAfter:  notBefore.Add(365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageDigitalSignature,
	}
	// Self-sign if no parent provided.
	if parent == nil {
		parent = template
		parentKey = key
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, &key.PublicKey, parentKey)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}
	certPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}))
	parsedCert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse generated certificate: %v", err)
	}
	return certPEM, parsedCert, key
}

func TestParseAndOrderCertificateChain(t *testing.T) {
	// Test case 1: empty chain.
	t.Run("empty chain", func(t *testing.T) {
		_, _, err := parseAndOrderCertificateChain("")
		if err == nil || !strings.Contains(err.Error(), "empty certificate chain") {
			t.Errorf("expected error for empty chain, got: %v", err)
		}
	})

	// Test case 2: no valid certificates found (wrong PEM type).
	t.Run("no valid certificates", func(t *testing.T) {
		invalid := "-----BEGIN NOT CERTIFICATE-----\nabc\n-----END NOT CERTIFICATE-----"
		_, _, err := parseAndOrderCertificateChain(invalid)
		if err == nil || !strings.Contains(err.Error(), "no valid certificates found") {
			t.Errorf("expected error for no valid certificates, got: %v", err)
		}
	})

	// Test case 3: single certificate (self-signed).
	t.Run("single certificate", func(t *testing.T) {
		notBefore := time.Now().Add(-1 * time.Hour)
		certPEM, _, _ := createCertificate(t, "single", "single", 100, notBefore, nil, nil)
		leaf, intermediates, err := parseAndOrderCertificateChain(certPEM)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Leaf should equal the original cert, with no intermediates.
		if string(leaf) != certPEM {
			t.Errorf("expected leaf to be the certificate, got different value")
		}
		if len(intermediates) != 0 {
			t.Errorf("expected no intermediates, got %d", len(intermediates))
		}
	})

	// Test case 4: valid chain with leaf, intermediate, and root in random order.
	t.Run("valid chain ordering", func(t *testing.T) {
		now := time.Now()
		// Create root certificate (self-signed).
		rootPEM, rootCert, rootKey := createCertificate(t, "root", "root", 1, now.Add(-10*time.Hour), nil, nil)
		// Create intermediate certificate signed by root.
		intermediatePEM, intermediateCert, intermediateKey := createCertificate(t, "intermediate", "root", 2, now.Add(-5*time.Hour), rootCert, rootKey)
		// Create leaf certificate signed by intermediate.
		leafPEM, _, _ := createCertificate(t, "leaf", "intermediate", 3, now.Add(-1*time.Hour), intermediateCert, intermediateKey)
		// Combine in random order: intermediate, leaf, root.
		chain := intermediatePEM + leafPEM + rootPEM
		leafOut, intermediates, err := parseAndOrderCertificateChain(chain)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Expected leaf is the leaf certificate.
		if string(leafOut) != leafPEM {
			t.Errorf("expected leaf to be leaf certificate")
		}
		// Expected intermediates: first should be the intermediate, second the root.
		if len(intermediates) < 2 {
			t.Errorf("expected at least 2 intermediates, got %d", len(intermediates))
		} else {
			if string(intermediates[0]) != intermediatePEM {
				t.Errorf("expected first intermediate to be intermediate certificate")
			}
			if string(intermediates[1]) != rootPEM {
				t.Errorf("expected second intermediate to be root certificate")
			}
		}
	})

	// Test case 5: valid chain with an extra certificate not connected to the chain.
	t.Run("chain with extra certificate", func(t *testing.T) {
		now := time.Now()
		// Create a proper chain: root, intermediate, leaf.
		rootPEM, rootCert, rootKey := createCertificate(t, "root", "root", 1, now.Add(-10*time.Hour), nil, nil)
		intermediatePEM, intermediateCert, intermediateKey := createCertificate(t, "intermediate", "root", 2, now.Add(-5*time.Hour), rootCert, rootKey)
		leafPEM, _, _ := createCertificate(t, "leaf", "intermediate", 3, now.Add(-1*time.Hour), intermediateCert, intermediateKey)
		// Create an extra self-signed certificate that doesn't chain.
		extraPEM, _, _ := createCertificate(t, "extra", "extra", 4, now.Add(-2*time.Hour), nil, nil)
		// Combine in random order: extra, root, leaf, intermediate.
		chain := extraPEM + rootPEM + leafPEM + intermediatePEM
		leafOut, intermediates, err := parseAndOrderCertificateChain(chain)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Expected leaf remains the leaf certificate.
		if string(leafOut) != leafPEM {
			t.Errorf("expected leaf to be leaf certificate")
		}
		// Intermediates should include the intermediate and root in order, with the extra appended.
		foundIntermediate, foundRoot, foundExtra := false, false, false
		for _, interm := range intermediates {
			if string(interm) == intermediatePEM {
				foundIntermediate = true
			}
			if string(interm) == rootPEM {
				foundRoot = true
			}
			if string(interm) == extraPEM {
				foundExtra = true
			}
		}
		if !foundIntermediate || !foundRoot || !foundExtra {
			t.Errorf("expected intermediates to contain intermediate, root, and extra certificates")
		}
	})

	// Test case 6: chain with an invalid certificate block.
	t.Run("invalid certificate block", func(t *testing.T) {
		invalidCert := "-----BEGIN CERTIFICATE-----\ninvalid\n-----END CERTIFICATE-----"
		chain := invalidCert
		_, _, err := parseAndOrderCertificateChain(chain)
		// Since no valid certificate is found, we expect the error to indicate that.
		if err == nil || !strings.Contains(err.Error(), "no valid certificates found") {
			t.Errorf("expected error for invalid certificate block, got: %v", err)
		}
	})

	// Test case 7: multiple self-signed certificates; select the one with the most recent NotBefore.
	t.Run("multiple self-signed certificates", func(t *testing.T) {
		now := time.Now()
		// Create two self-signed certificates with different NotBefore times.
		certAPEM, _, _ := createCertificate(t, "A", "A", 10, now.Add(-2*time.Hour), nil, nil)
		certBPEM, _, _ := createCertificate(t, "B", "B", 11, now.Add(-1*time.Hour), nil, nil)
		// Combine in any order.
		chain := certAPEM + certBPEM
		leafOut, intermediates, err := parseAndOrderCertificateChain(chain)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Expect leaf to be certificate B since its NotBefore is later.
		if string(leafOut) != certBPEM {
			t.Errorf("expected leaf to be certificate B")
		}
		// The other certificate should appear as an intermediate.
		if len(intermediates) != 1 || string(intermediates[0]) != certAPEM {
			t.Errorf("expected certificate A to be intermediate")
		}
	})
}
