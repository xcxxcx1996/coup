package model

import (
	"errors"
)

func (s *MatchState) GainCoins(playerId string, coins int32) error {
	player, ok := s.PlayerInfos[playerId]
	if !ok {
		return errors.New("no record")
	}
	player.Coins += (coins)
	return nil
}
func (s *MatchState) LoseCoins(playerId string, coins int32) error {
	player, ok := s.PlayerInfos[playerId]
	if !ok {
		return errors.New("no record")
	}
	player.Coins -= (coins)
	return nil
}

func (s *MatchState) SetCoins(playerId string, coins int32) error {
	player, ok := s.PlayerInfos[playerId]
	if !ok {
		return errors.New("no record")
	}
	player.Coins = (coins)
	return nil
}

func (s *MatchState) GetCoins(playerId string) (int32, error) {
	player, ok := s.PlayerInfos[playerId]
	if !ok {
		return 0, errors.New("no record")
	}
	return player.Coins, nil
}
