package archivista

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// parseAndOrderCertificateChain parses a PEM-encoded certificate chain string,
// validates each "CERTIFICATE" PEM block, and orders them so that the leaf certificate is first,
// followed by intermediates. If no intermediates are present, it simply returns the single certificate.
func parseAndOrderCertificateChain(chain string) (leafCert []byte, intermediates [][]byte, err error) {
	if chain == "" {
		return nil, nil, fmt.Errorf("empty certificate chain")
	}

	data := []byte(chain)
	type parsedCert struct {
		cert *x509.Certificate
		pem  []byte
	}
	var certs []parsedCert

	// Parse all PEM blocks of type "CERTIFICATE"
	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			continue
		}
		parsed, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse certificate: %w", err)
		}
		certs = append(certs, parsedCert{
			cert: parsed,
			pem:  pem.EncodeToMemory(block),
		})
	}

	if len(certs) == 0 {
		return nil, nil, fmt.Errorf("no valid certificates found in chain")
	}

	// If only one certificate is found, return it as the leaf.
	if len(certs) == 1 {
		return certs[0].pem, nil, nil
	}

	// Identify leaf candidates:
	// Prefer certificates that are not self-signed and whose Subject isn't used as an Issuer in any other cert.
	var leafCandidates []parsedCert
	for i, pc := range certs {
		selfSigned := pc.cert.Subject.String() == pc.cert.Issuer.String()
		if !selfSigned {
			used := false
			for j, other := range certs {
				if i == j {
					continue
				}
				if other.cert.Issuer.String() == pc.cert.Subject.String() {
					used = true
					break
				}
			}
			if !used {
				leafCandidates = append(leafCandidates, pc)
			}
		}
	}
	// If no non-self-signed candidate, fall back to self-signed ones.
	if len(leafCandidates) == 0 {
		leafCandidates = certs
	}

	// Choose the best leaf candidate: if multiple, select the one with the most recent NotBefore date.
	leaf := leafCandidates[0]
	for _, candidate := range leafCandidates[1:] {
		if candidate.cert.NotBefore.After(leaf.cert.NotBefore) {
			leaf = candidate
		}
	}

	// Build a map for quick lookup (subject => parsedCert).
	subjectMap := make(map[string]parsedCert)
	for _, pc := range certs {
		subjectMap[pc.cert.Subject.String()] = pc
	}

	// Order the chain starting from the leaf.
	ordered := []parsedCert{leaf}
	used := map[string]bool{leaf.cert.SerialNumber.String(): true}
	current := leaf
	for {
		next, found := subjectMap[current.cert.Issuer.String()]
		if !found || used[next.cert.SerialNumber.String()] {
			break
		}
		ordered = append(ordered, next)
		used[next.cert.SerialNumber.String()] = true
		current = next
	}

	leafCert = ordered[0].pem
	// The intermediates are any ordered certs after the leaf.
	for i := 1; i < len(ordered); i++ {
		intermediates = append(intermediates, ordered[i].pem)
	}
	// Append any extra certificates not included in the ordering.
	for _, pc := range certs {
		if !used[pc.cert.SerialNumber.String()] {
			intermediates = append(intermediates, pc.pem)
		}
	}

	return leafCert, intermediates, nil
}
