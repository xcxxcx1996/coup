package state

func (s *MatchState) ResetDeadLine() {
	s.DeadlineRemainingTicks = s.MaxDeadlineTicks
}

func (s *MatchState) ResetNextMartch() {
	s.NextGameRemainingTicks = 40000
}
