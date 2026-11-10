package env

// Variables represents a collection of environment variables
// as a map of string keys to string values.
type Variables map[string]string

// MustGet retrieves the value of the environment variable with the given key
// and panics if the variable is not found.
func (v Variables) MustGet(k string) string {
	if value, ok := v.Get(k); ok {
		return value
	}
	panic("variable not found: " + k)
}

// Get retrieves the value of the environment variable with the given key.
//
// It returns the value and a boolean indicating whether the variable was found.
func (v Variables) Get(k string) (string, bool) {
	if value, ok := v[k]; ok {
		return value, true
	}
	return "", false
}

// Set sets the value of the environment variable
// with the given key to the specified value.
func (v Variables) Set(k, value string) {
	v[k] = value
}
