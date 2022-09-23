package object

import (
	"github.com/icheka/sonar-lang/sonar-lang/errors"
)

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{Store: s, outer: nil, Readonly: make(map[string]bool)}
}

func NewEphemeralScope(allow []string, readonly map[string]bool, parent *Environment) *Environment {
	env := NewEnvironment()
	env.allow = allow
	env.Readonly = readonly
	env.outer = parent
	return env
}

type Environment struct {
	Store    map[string]Object
	outer    *Environment
	allow    []string
	Readonly map[string]bool
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.Store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	// if this variable is a constant, it is immutable
	// - check whether it's been initialised as a constant,
	// - if it has, return an error
	if e.Store[name] != nil && e.Readonly[name] {
		return &Error{Conf: errors.ConstantAssignmentError(name)}
	}

	if len(e.allow) > 0 {
		for _, identifier := range e.allow {
			if identifier == name {
				e.Store[name] = val
				return val
			}
		}

		// `name` cannot be stored in e.Store, try storing in e.outer.Store instead
		if e.outer != nil {
			return e.outer.Set(name, val)
		}
	}
	e.Store[name] = val
	return val
}
