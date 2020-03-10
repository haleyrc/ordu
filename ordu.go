package ordu

import (
	"context"
	"fmt"
	"os"
)

type ExecutionContext struct {
	Ctx         context.Context
	Args        []string
	Environment Environment
}

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

type Runner interface {
	Run(ec ExecutionContext) error
}

type Dispatch map[string]Runner

func NewManager() (Manager, error) {
	ec, err := NewExecutionContext()
	if err != nil {
		return Manager{}, fmt.Errorf("NewManager: %w", err)
	}
	return Manager{ec: ec}, nil
}

type Manager struct {
	ec       ExecutionContext
	Dispatch Dispatch
}

func (m Manager) PrintCommands() {
	for cmd := range m.Dispatch {
		fmt.Fprintln(os.Stderr, cmd)
	}
}

func (m Manager) Run(cmd string) error {
	runner, ok := m.Dispatch[cmd]
	if !ok {
		return fmt.Errorf("Manager.Run: invalid command: %s", cmd)
	}
	return runner.Run(ExecutionContext{
		Ctx:         m.ec.Ctx,
		Environment: m.ec.Environment,
		Args:        m.ec.Args[2:],
	})
}
