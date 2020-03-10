package ordu

import (
	"errors"
	"os"
	"strings"
)

type Environment map[string]string

func (e Environment) Get(key string) string {
	val, _ := e.Lookup(key)
	return val
}

func (e Environment) Lookup(key string) (string, bool) {
	val, ok := e[key]
	return val, ok
}

func LoadEnvironment() (Environment, error) {
	return loadEnvironment(make(Environment))
}

func loadEnvironment(e Environment) (Environment, error) {
	all := os.Environ()
	for _, one := range all {
		parts := strings.SplitN(one, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("ordu: LoadEnvironment: invalid number of segments")
		}
		e[parts[0]] = parts[1]
	}
	return e, nil
}

func LoadEnvironmentWithDefaults(def Environment) (Environment, error) {
	return loadEnvironment(def)
}
