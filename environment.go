package ordu

import (
	"errors"
	"os"
	"strings"
)

// Environment is a key:value representation of a set of environment variables.
type Environment map[string]string

// Get returns the value of the environment variable identified by key if it
// exists, or a blank string otherwise.
func (e Environment) Get(key string) string {
	val, _ := e.Lookup(key)
	return val
}

// Lookup returns the value of the environment variable identified by key. If no
// variable with that name is found, it will return a blank string and false.
func (e Environment) Lookup(key string) (string, bool) {
	val, ok := e[key]
	return val, ok
}

// LoadEnvironment returns an Environment populated by the variables reported by
// the OS.
func LoadEnvironment() (Environment, error) {
	return loadEnvironment(make(Environment))
}

// LoadEnvironmentWithDefaults returns an Environment populated by the variables
// reported by the OS. Any defaults provided will be overwritten by the OS
// reported value if any, otherwise the default will be available.
func LoadEnvironmentWithDefaults(def Environment) (Environment, error) {
	return loadEnvironment(def)
}

func loadEnvironment(e Environment) (Environment, error) {
	all := os.Environ()
	for _, one := range all {
		parts := strings.SplitN(one, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("LoadEnvironment: invalid number of segments")
		}
		e[parts[0]] = parts[1]
	}
	return e, nil
}
