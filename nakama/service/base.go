package service

type MatchService struct {
	// LastAction int
	// Time       int
}

func New() (s *MatchService) {
	return new(MatchService)
}
