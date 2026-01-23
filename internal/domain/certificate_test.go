package domain

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCertificate_Success(t *testing.T) {
	// Create a temporary valid PEM file
	tmpDir := t.TempDir()
	certPath := filepath.Join(tmpDir, "test-ca.pem")

	validPEM := `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHBfpegPjMCMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RjYTAeFw0yNDAxMDEwMDAwMDBaFw0yNTAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RjYTBcMA0GCSqGSIb3DQEBAQUAA0sAMEgCQQC7o96FCpCKgJ0P8QFxYVJr
4PXZhCBDWNdkG/V7ZlFNL/JVFpN0hG9VKPLUPkTjELew6b1AKsRzwKhzqZjTwEzl
AgMBAAGjUzBRMB0GA1UdDgQWBBQTRYmVpPHvfmVDPFMqdCb9IYbffjAfBgNVHSME
GDAWgBQTRYmVpPHvfmVDPFMqdCb9IYbffjAPBgNVHRMBAf8EBTADAQH/MA0GCSqG
SIb3DQEBCwUAA0EAjR3OWp5k8PVZ3zRV8hB4IU7DaGYXAiH3b8PxUXGpfHxvdT0Q
C/3g2h9qB4mH0JOz4vCnLNL5k3x5M2l1AJVNpA==
-----END CERTIFICATE-----`

	err := os.WriteFile(certPath, []byte(validPEM), 0644)
	if err != nil {
		t.Fatalf("Failed to write test cert: %v", err)
	}

	cert, err := LoadCertificate(certPath)
	if err != nil {
		t.Fatalf("LoadCertificate() error = %v, want nil", err)
	}

	if cert.Path != certPath {
		t.Errorf("Certificate.Path = %v, want %v", cert.Path, certPath)
	}
	if len(cert.Content) == 0 {
		t.Error("Certificate.Content is empty")
	}
}

func TestLoadCertificate_FileNotFound(t *testing.T) {
	_, err := LoadCertificate("/nonexistent/path/to/cert.pem")

	if err == nil {
		t.Fatal("LoadCertificate() error = nil, want error")
	}
	if !errors.Is(err, ErrCertificateNotFound) {
		t.Errorf("LoadCertificate() error = %v, want %v", err, ErrCertificateNotFound)
	}
}

func TestLoadCertificate_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	certPath := filepath.Join(tmpDir, "empty.pem")

	err := os.WriteFile(certPath, []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to write empty file: %v", err)
	}

	_, err = LoadCertificate(certPath)

	if err == nil {
		t.Fatal("LoadCertificate() error = nil, want error")
	}
	if !errors.Is(err, ErrCertificateEmpty) {
		t.Errorf("LoadCertificate() error = %v, want %v", err, ErrCertificateEmpty)
	}
}

func TestLoadCertificate_InvalidPEM(t *testing.T) {
	tmpDir := t.TempDir()
	certPath := filepath.Join(tmpDir, "invalid.pem")

	err := os.WriteFile(certPath, []byte("not a valid PEM file"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid file: %v", err)
	}

	_, err = LoadCertificate(certPath)

	if err == nil {
		t.Fatal("LoadCertificate() error = nil, want error")
	}
	if !errors.Is(err, ErrCertificateInvalidPEM) {
		t.Errorf("LoadCertificate() error = %v, want %v", err, ErrCertificateInvalidPEM)
	}
}

func TestLoadCertificate_NotCertificate(t *testing.T) {
	tmpDir := t.TempDir()
	certPath := filepath.Join(tmpDir, "key.pem")

	// This is a PEM but it's a private key, not a certificate
	keyPEM := `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBALuj3oUKkIqAnQ/xAXFhUmvg9dmEIENY12Qb9XtmUU0v8lUWk3SE
b1Uo8tQ+ROMQ
-----END RSA PRIVATE KEY-----`

	err := os.WriteFile(certPath, []byte(keyPEM), 0644)
	if err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	_, err = LoadCertificate(certPath)

	if err == nil {
		t.Fatal("LoadCertificate() error = nil, want error")
	}
	if !errors.Is(err, ErrCertificateNotCA) {
		t.Errorf("LoadCertificate() error = %v, want %v", err, ErrCertificateNotCA)
	}
}

func TestCertificate_String(t *testing.T) {
	cert := Certificate{
		Path:    "/path/to/ca.pem",
		Content: []byte("test content"),
	}

	got := cert.String()
	expected := "Certificate(/path/to/ca.pem, 12 bytes)"

	if got != expected {
		t.Errorf("Certificate.String() = %v, want %v", got, expected)
	}
}

func TestCertificate_PEMString(t *testing.T) {
	content := "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----"
	cert := Certificate{
		Path:    "/path/to/ca.pem",
		Content: []byte(content),
	}

	if got := cert.PEMString(); got != content {
		t.Errorf("Certificate.PEMString() = %v, want %v", got, content)
	}
}
