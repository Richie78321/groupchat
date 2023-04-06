package util

type Signal struct {
	// A channel for signalling. Needs a strict buffer size of 1.
	signal chan struct{}
}

func NewSignal() *Signal {
	return &Signal{
		// The buffer size is strictly 1 to ensure proper signalling behavior.
		signal: make(chan struct{}, 1),
	}
}

func (s *Signal) Signal() {
	select {
	case s.signal <- struct{}{}:
	default:
		// The signal channel has a strict buffer size of 1.
		// If there is already a signal in the channel,
		// then we can safely continue because the channel has already
		// been signalled.
	}
}

func (s *Signal) GetSignal() <-chan struct{} {
	return s.signal
}
