package port

import "context"

// System defines the contract for interacting with the host operating system.
//
// This port abstracts all OS-level operations:
//   - Filesystem access
//   - Environment variables
//   - Process execution
//
// The Core NEVER depends on this interface directly.
// Only Adapters use System to interact with the host.
//
// This abstraction enables:
//   - Unit testing with mock implementations
//   - Cross-platform compatibility
//   - Controlled side effects
type System interface {
	// --- Filesystem Operations ---

	// FileExists checks if a file exists at the given path.
	FileExists(path string) bool

	// DirExists checks if a directory exists at the given path.
	DirExists(path string) bool

	// ReadFile reads the entire contents of a file.
	ReadFile(path string) ([]byte, error)

	// WriteFile writes data to a file, creating it if necessary.
	// Overwrites existing content.
	WriteFile(path string, data []byte, perm uint32) error

	// AppendToFile appends data to a file, creating it if necessary.
	AppendToFile(path string, data []byte) error

	// MkdirAll creates a directory and all parent directories.
	MkdirAll(path string, perm uint32) error

	// --- Environment Variables ---

	// GetEnv returns the value of an environment variable.
	// Returns empty string if not set.
	GetEnv(key string) string

	// SetEnv sets an environment variable for the current process.
	SetEnv(key, value string) error

	// --- Process Execution ---

	// RunCommand executes a command and returns its output.
	// Returns stdout, combined with stderr on error.
	RunCommand(ctx context.Context, name string, args ...string) ([]byte, error)

	// RunCommandSilent executes a command and returns only the exit status.
	// Output is discarded. Useful for detection.
	RunCommandSilent(ctx context.Context, name string, args ...string) error

	// CommandExists checks if a command is available in PATH.
	CommandExists(name string) bool

	// --- Paths ---

	// HomeDir returns the current user's home directory.
	HomeDir() (string, error)

	// TempDir returns the system's temporary directory.
	TempDir() string
}
