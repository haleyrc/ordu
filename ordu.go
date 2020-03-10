package ordu

import (
	"context"
	"fmt"
	"io"
	"os"
)

// ExecutionContext contains all of the runtime configuration that gets passed
// to commands.
type ExecutionContext struct {
	// Ctx is originally set to context.Background(), but can be overridden if
	// required.
	Ctx context.Context

	// Args is originally set to all of the command line arguments, but only
	// those following the command name are passed to individual commands.
	Args []string

	// Environment is a key:value store of all of the environment variables
	// reported by the operating system.
	Environment Environment
}

// NewExecutionContext returns an ExecutionContext suitable for passing to
// commands.
func NewExecutionContext() (ExecutionContext, error) {
	env, err := LoadEnvironment()
	if err != nil {
		return ExecutionContext{}, fmt.Errorf("NewExecutionContext: %w", err)
	}
	return ExecutionContext{
		Ctx:         context.Background(),
		Args:        os.Args,
		Environment: env,
	}, nil
}

// Runner must be implemented by any commands.
type Runner interface {
	// Run is called by a Manager if a command's key matches the command passed
	// to the application.
	Run(ec ExecutionContext) error
}

// Dispatch is a mapping from a command name, which is taken as an argument on
// the command line, to a Runner.
type Dispatch map[string]Runner

// NewManager returns a manager for running commands.
func NewManager() (Manager, error) {
	ec, err := NewExecutionContext()
	if err != nil {
		return Manager{}, fmt.Errorf("NewManager: %w", err)
	}
	return Manager{
		ec:       ec,
		Messages: os.Stderr,
	}, nil
}

// Manager provides the ability to run commands by providing a command name on
// the command line. Manager is responsible for passing the ExecutionContext it
// was created with (after stripping the command name from args), to the command
// that is being invoked.
type Manager struct {
	ec       ExecutionContext
	Dispatch Dispatch
	Messages io.Writer
}

// PrintCommands prints all of the commands registered in the Manager's Dispatch
// field.
func (m Manager) PrintCommands() {
	for cmd := range m.Dispatch {
		fmt.Fprintln(m.Messages, cmd)
	}
}

// Run attempts to run the command identified by name. It returns an error if a
// command with that name does not exist, or if the command invoked returns an
// error.
func (m Manager) Run(name string) error {
	runner, ok := m.Dispatch[name]
	if !ok {
		return fmt.Errorf("Manager.Run: invalid command: %s", name)
	}
	return runner.Run(ExecutionContext{
		Ctx:         m.ec.Ctx,
		Environment: m.ec.Environment,
		Args:        m.ec.Args[2:],
	})
}
