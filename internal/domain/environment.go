package domain

// Environment represents a logical trust domain.
//
// An Environment is NOT:
//   - A binary or executable
//   - A command or process
//   - A package manager
//
// An Environment IS:
//   - A conceptual boundary where trust rules apply
//   - A runtime or tool ecosystem with its own certificate trust model
//
// Examples: OpenSSL, Python (certifi), Node.js, Git, JVM
//
// Environment is a value object - it identifies an environment but contains
// no behavior. The behavior lives in the adapters that implement EnvironmentPort.
type Environment struct {
	// ID is a unique identifier for this environment.
	// Used for programmatic access and configuration.
	// Examples: "openssl", "python", "nodejs", "git"
	ID string

	// Name is a human-readable display name.
	// Used for CLI output and logging.
	// Examples: "OpenSSL", "Python (certifi)", "Node.js", "Git"
	Name string
}

// NewEnvironment creates a new Environment with the given ID and name.
func NewEnvironment(id, name string) Environment {
	return Environment{
		ID:   id,
		Name: name,
	}
}

// String returns the human-readable name of the environment.
func (e Environment) String() string {
	return e.Name
}
