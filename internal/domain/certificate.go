package domain

import (
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Certificate errors.
var (
	// ErrCertificateNotFound is returned when the certificate file does not exist.
	ErrCertificateNotFound = errors.New("certificate file not found")

	// ErrCertificateEmpty is returned when the certificate file is empty.
	ErrCertificateEmpty = errors.New("certificate file is empty")

	// ErrCertificateInvalidPEM is returned when the file is not valid PEM format.
	ErrCertificateInvalidPEM = errors.New("certificate is not valid PEM format")

	// ErrCertificateNotCA is returned when the certificate is not a CA certificate.
	ErrCertificateNotCA = errors.New("certificate does not appear to be a CA certificate")
)

// Certificate represents a CA certificate to be trusted.
//
// A Certificate is loaded from a PEM file and validated on construction.
// This ensures we fail fast if the certificate is invalid, rather than
// discovering problems during the Apply phase.
//
// Certificate is a value object - immutable after creation.
type Certificate struct {
	// Path is the filesystem path where the certificate was loaded from.
	Path string

	// Content is the raw PEM-encoded certificate data.
	Content []byte
}

// LoadCertificate loads and validates a CA certificate from the given path.
//
// Validation includes:
//   - File exists and is readable
//   - File is not empty
//   - Content is valid PEM format
//   - PEM block type suggests a certificate (CERTIFICATE or X509 CERTIFICATE)
//
// Returns ErrCertificateNotFound, ErrCertificateEmpty, ErrCertificateInvalidPEM,
// or ErrCertificateNotCA on validation failure.
func LoadCertificate(path string) (Certificate, error) {
	// Check file exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return Certificate{}, fmt.Errorf("%w: %s", ErrCertificateNotFound, path)
	}
	if err != nil {
		return Certificate{}, fmt.Errorf("failed to stat certificate file: %w", err)
	}

	// Check not empty
	if info.Size() == 0 {
		return Certificate{}, fmt.Errorf("%w: %s", ErrCertificateEmpty, path)
	}

	// Read content
	content, err := os.ReadFile(path)
	if err != nil {
		return Certificate{}, fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Validate PEM format
	block, _ := pem.Decode(content)
	if block == nil {
		return Certificate{}, fmt.Errorf("%w: %s", ErrCertificateInvalidPEM, path)
	}

	// Check it looks like a certificate (not a key or other PEM type)
	blockType := strings.ToUpper(block.Type)
	if !strings.Contains(blockType, "CERTIFICATE") {
		return Certificate{}, fmt.Errorf("%w: PEM type is %q, expected CERTIFICATE", ErrCertificateNotCA, block.Type)
	}

	return Certificate{
		Path:    path,
		Content: content,
	}, nil
}

// String returns a short description of the certificate.
func (c Certificate) String() string {
	return fmt.Sprintf("Certificate(%s, %d bytes)", c.Path, len(c.Content))
}

// PEMString returns the certificate content as a string.
// Useful for appending to bundle files.
func (c Certificate) PEMString() string {
	return string(c.Content)
}
