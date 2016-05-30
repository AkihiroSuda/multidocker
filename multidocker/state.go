package multidocker

// State represents the state of a host
type State int

const (
	// UnknownState is an unknown state.
	UnknownState State = iota
	// Running denotes that the host is running.
	Running
	// Stopped denotes that the host is not running.
	Stopped
)

func (s State) String() string {
	switch s {
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	default:
		return "UnknownState?"
	}
}
