package system

import (
	"context"
	"os"
	"os/exec"
)

// OSSystem is the real implementation of port.System.
// It delegates all operations to the actual operating system.
type OSSystem struct{}

// New creates a new OSSystem.
func New() *OSSystem {
	return &OSSystem{}
}

// --- Filesystem Operations ---

// FileExists checks if a file exists at the given path.
func (s *OSSystem) FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists at the given path.
func (s *OSSystem) DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ReadFile reads the entire contents of a file.
func (s *OSSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes data to a file, creating it if necessary.
func (s *OSSystem) WriteFile(path string, data []byte, perm uint32) error {
	return os.WriteFile(path, data, os.FileMode(perm))
}

// AppendToFile appends data to a file, creating it if necessary.
func (s *OSSystem) AppendToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

// MkdirAll creates a directory and all parent directories.
func (s *OSSystem) MkdirAll(path string, perm uint32) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

// --- Environment Variables ---

// GetEnv returns the value of an environment variable.
func (s *OSSystem) GetEnv(key string) string {
	return os.Getenv(key)
}

// SetEnv sets an environment variable for the current process.
func (s *OSSystem) SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// --- Process Execution ---

// RunCommand executes a command and returns its output.
func (s *OSSystem) RunCommand(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.CombinedOutput()
}

// RunCommandSilent executes a command and returns only the exit status.
func (s *OSSystem) RunCommandSilent(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// CommandExists checks if a command is available in PATH.
func (s *OSSystem) CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// --- Paths ---

// HomeDir returns the current user's home directory.
func (s *OSSystem) HomeDir() (string, error) {
	return os.UserHomeDir()
}

// TempDir returns the system's temporary directory.
func (s *OSSystem) TempDir() string {
	return os.TempDir()
}
