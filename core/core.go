package core

// Instrumenter is an interface for instrumenting a core.
type Instrumenter interface {
	Instrument(string) ([]byte, error)
}
