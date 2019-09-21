package vm

import "time"

type Stats struct {
	InitAt, StartedAt, StoppedAt, PausedAt, SuspendedAt time.Time
	Instructions                                        int
	Defers                                              int
}

func (s *Stats) recordTime(oldState, newState State) {
	if oldState == Init {
		s.StartedAt = time.Now()
	}

	switch newState {
	case Stopped:
		s.StoppedAt = time.Now()
	case Paused:
		s.PausedAt = time.Now()
	case Suspended:
		s.SuspendedAt = time.Now()
	}
}
