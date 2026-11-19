package connector

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"

	googlecredentialguard "github.com/conductorone/google-credential-guard"
)

func writeCredentials(t *testing.T, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "credentials.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func serviceAccountJSON(t *testing.T) map[string]string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return map[string]string{
		"type":                        "service_account",
		"project_id":                  "test-project",
		"private_key_id":              "test-key",
		"private_key":                 string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKey})),
		"client_email":                "test@test-project.iam.gserviceaccount.com",
		"client_id":                   "123456789012345678901",
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/test",
		"token_uri":                   googlecredentialguard.TokenURL,
	}
}

func TestCredentialsOptionAcceptsConsoleKey(t *testing.T) {
	data, err := json.Marshal(serviceAccountJSON(t))
	if err != nil {
		t.Fatal(err)
	}
	path := writeCredentials(t, data)
	if _, err := credentialsOption(path); err != nil {
		t.Fatalf("valid service account rejected: %v", err)
	}
	// A malformed FIREBASE_CONFIG path must not make Firebase read another file.
	t.Setenv("FIREBASE_CONFIG", filepath.Join(t.TempDir(), "missing"))
	if _, err := New(context.Background(), path); err != nil {
		t.Fatalf("valid service account rejected by Firebase: %v", err)
	}
}

func TestNewRejectsUnsafeCredentialFiles(t *testing.T) {
	valid := serviceAccountJSON(t)
	valid["token_uri"] = "http://127.0.0.1/token"
	hostileTokenURI, err := json.Marshal(valid)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"external account file", []byte(`{"type":"external_account","credential_source":{"file":"/tmp/subject-token"}}`)},
		{"external account URL", []byte(`{"type":"external_account","credential_source":{"url":"https://example.invalid/token"}}`)},
		{"external account executable", []byte(`{"type":"external_account","credential_source":{"executable":{"command":"id"}}}`)},
		{"hostile token URI", hostileTokenURI},
		{"oversized", []byte(strings.Repeat("x", googlecredentialguard.MaxJSONBytes+1))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(context.Background(), writeCredentials(t, tc.data))
			if err == nil || !strings.Contains(err.Error(), "validate credentials file") {
				t.Fatalf("expected validation error, got %v", err)
			}
			if strings.Contains(err.Error(), "127.0.0.1") || strings.Contains(err.Error(), "/tmp/subject-token") {
				t.Fatalf("credential content leaked in error: %v", err)
			}
		})
	}
}

func TestCredentialsOptionSanitizesPathError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "private-file-name.json")
	_, err := credentialsOption(path)
	if err == nil || strings.Contains(err.Error(), path) || strings.Contains(err.Error(), "private-file-name") {
		t.Fatalf("path leaked in error: %v", err)
	}
}
