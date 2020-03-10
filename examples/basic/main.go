package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/haleyrc/ordu"
)

func main() {
	mgr, err := ordu.NewManager()
	if err != nil {
		panic(err)
	}
	mgr.Dispatch = ordu.Dispatch{
		"greet":  Greeter("Hello"),
		"mangle": NewMangler(),
		"lang":   LanguageReporter("en_US.UTF-8"),
	}
	if len(os.Args) == 1 {
		mgr.PrintCommands()
		os.Exit(1)
	}
	if err := mgr.Run(os.Args[1]); err != nil {
		panic(err)
	}
}

type Greeter string

func (g Greeter) Run(ec ordu.ExecutionContext) error {
	fs := flag.NewFlagSet("greeter", flag.PanicOnError)
	name := fs.String("name", "World", "The name of the person to greet")
	fs.Parse(ec.Args)
	fmt.Printf("%s, %s!\n", g, *name)
	return nil
}

func NewMangler() Mangler {
	return Mangler{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type Mangler struct {
	r *rand.Rand
}

func (m Mangler) Run(ec ordu.ExecutionContext) error {
	if len(ec.Args) < 1 {
		return errors.New("Mangler.Run: must provide a string to mangle")
	}
	_, err := fmt.Println(m.mangleString(ec.Args[0]))
	if err != nil {
		return fmt.Errorf("Mangler.Run: %w", err)
	}
	return nil
}

func (m Mangler) mangleString(s string) string {
	var sb strings.Builder
	for _, c := range s {
		sb.WriteRune(m.mangleByte(c))
	}
	return sb.String()
}

func (m Mangler) mangleByte(r rune) rune {
	if f := m.r.Float32(); f < .5 {
		return r
	}
	if r >= 65 && r <= 90 {
		// lowercase it
		return r + 32
	}
	if r >= 97 || r <= 122 {
		// uppercase it
		return r - 32
	}
	return r
}

type LanguageReporter string

func (lr LanguageReporter) Run(ec ordu.ExecutionContext) error {
	lang, ok := ec.Environment.Lookup("LANG")
	if !ok {
		lang = string(lr)
	}
	fmt.Println(lang)
	return nil
}
