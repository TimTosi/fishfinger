package fishfinger

// -----------------------------------------------------------------------------

// Origin is an `interface` used to manage a Docker environment
// programmatically.
type Origin interface {
	Start(...string) error
	StartBackoff(func(string) error, ...string) error
	Stop(...string) error
	Status(string) (bool, error)
	Port(string, string) (string, error)
	Env(string, string)
}
