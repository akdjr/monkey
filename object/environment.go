package object

// Environment represents an object environment.  This is where we keep track of all identifiers and their values
type Environment struct {
	store map[string]Object
}

// NewEnvironment creates an empty environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Get returns the Object stored at name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set creates an entry for Object at name
func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
