package system

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	sys := New()
	if sys == nil {
		t.Fatal("New() returned nil")
	}
}

func TestOSSystem_FileExists(t *testing.T) {
	sys := New()

	// Create temp file
	f, err := os.CreateTemp("", "trustica-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	if !sys.FileExists(f.Name()) {
		t.Error("FileExists() should return true for existing file")
	}

	if sys.FileExists("/nonexistent/path/file.txt") {
		t.Error("FileExists() should return false for non-existing file")
	}

	// Directory should not count as file
	if sys.FileExists(os.TempDir()) {
		t.Error("FileExists() should return false for directory")
	}
}

func TestOSSystem_DirExists(t *testing.T) {
	sys := New()

	if !sys.DirExists(os.TempDir()) {
		t.Error("DirExists() should return true for existing directory")
	}

	if sys.DirExists("/nonexistent/path/dir") {
		t.Error("DirExists() should return false for non-existing directory")
	}

	// File should not count as directory
	f, err := os.CreateTemp("", "trustica-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	if sys.DirExists(f.Name()) {
		t.Error("DirExists() should return false for file")
	}
}

func TestOSSystem_ReadWriteFile(t *testing.T) {
	sys := New()

	// Create temp dir
	dir, err := os.MkdirTemp("", "trustica-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "test.txt")
	content := []byte("hello, trustica!")

	// Write
	if err := sys.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	// Read
	data, err := sys.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}

	if string(data) != string(content) {
		t.Errorf("ReadFile() = %q, want %q", data, content)
	}
}

func TestOSSystem_AppendToFile(t *testing.T) {
	sys := New()

	dir, err := os.MkdirTemp("", "trustica-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "append.txt")

	// First append (creates file)
	if err := sys.AppendToFile(path, []byte("line1\n")); err != nil {
		t.Fatalf("First AppendToFile() error: %v", err)
	}

	// Second append
	if err := sys.AppendToFile(path, []byte("line2\n")); err != nil {
		t.Fatalf("Second AppendToFile() error: %v", err)
	}

	// Verify
	data, _ := sys.ReadFile(path)
	expected := "line1\nline2\n"
	if string(data) != expected {
		t.Errorf("AppendToFile result = %q, want %q", data, expected)
	}
}

func TestOSSystem_MkdirAll(t *testing.T) {
	sys := New()

	dir, err := os.MkdirTemp("", "trustica-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	nested := filepath.Join(dir, "a", "b", "c")

	if err := sys.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("MkdirAll() error: %v", err)
	}

	if !sys.DirExists(nested) {
		t.Error("MkdirAll() did not create nested directories")
	}
}

func TestOSSystem_GetSetEnv(t *testing.T) {
	sys := New()

	key := "TRUSTICA_TEST_VAR"
	value := "test_value_123"

	// Clean up after test
	original := os.Getenv(key)
	defer os.Setenv(key, original)

	if err := sys.SetEnv(key, value); err != nil {
		t.Fatalf("SetEnv() error: %v", err)
	}

	got := sys.GetEnv(key)
	if got != value {
		t.Errorf("GetEnv() = %q, want %q", got, value)
	}

	// Non-existent should return empty
	if sys.GetEnv("TRUSTICA_NONEXISTENT_VAR_XYZ") != "" {
		t.Error("GetEnv() should return empty for non-existent var")
	}
}

func TestOSSystem_RunCommand(t *testing.T) {
	sys := New()
	ctx := context.Background()

	// Simple echo command
	output, err := sys.RunCommand(ctx, "echo", "hello")
	if err != nil {
		t.Fatalf("RunCommand() error: %v", err)
	}

	// Output includes newline
	if string(output) != "hello\n" {
		t.Errorf("RunCommand() output = %q, want %q", output, "hello\n")
	}
}

func TestOSSystem_RunCommandSilent(t *testing.T) {
	sys := New()
	ctx := context.Background()

	// Should succeed silently
	if err := sys.RunCommandSilent(ctx, "true"); err != nil {
		t.Errorf("RunCommandSilent(true) should not error: %v", err)
	}

	// Should fail silently
	if err := sys.RunCommandSilent(ctx, "false"); err == nil {
		t.Error("RunCommandSilent(false) should error")
	}
}

func TestOSSystem_CommandExists(t *testing.T) {
	sys := New()

	// These should exist on any Unix system
	if !sys.CommandExists("echo") {
		t.Error("CommandExists(echo) should be true")
	}

	if sys.CommandExists("nonexistent_command_xyz123") {
		t.Error("CommandExists(nonexistent) should be false")
	}
}

func TestOSSystem_HomeDir(t *testing.T) {
	sys := New()

	home, err := sys.HomeDir()
	if err != nil {
		t.Fatalf("HomeDir() error: %v", err)
	}

	if home == "" {
		t.Error("HomeDir() returned empty string")
	}

	if !sys.DirExists(home) {
		t.Error("HomeDir() returned non-existent directory")
	}
}

func TestOSSystem_TempDir(t *testing.T) {
	sys := New()

	tmp := sys.TempDir()
	if tmp == "" {
		t.Error("TempDir() returned empty string")
	}

	if !sys.DirExists(tmp) {
		t.Error("TempDir() returned non-existent directory")
	}
}
