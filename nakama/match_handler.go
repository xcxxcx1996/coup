package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/service"
	model "github.com/xcxcx1996/coup/state"
	"google.golang.org/protobuf/proto"
)

const (
	moduleName = "coup"

	tickRate    = 5
	maxEmptySec = 5
	// delayBetweenGamesSec = 5
	turnTimeFastSec   = 10
	nextStartSec      = 5
	turnTimeNormalSec = 30
)

// Compile-time check to make sure all required functions are implemented.
var serv = service.New()
var _ runtime.Match = &MatchHandler{
	service: serv,
}

type MatchHandler struct {
	service *service.MatchService
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	fast, ok := params["fast"].(bool)
	if !ok {
		logger.Error("invalid match init parameter \"fast\"")
		return nil, 0, ""
	}
	size, ok := params["size"].(int64)
	if !ok {
		logger.Error("invalid match init parameter \"fast\"")
		return nil, 0, ""
	}
	Label := &model.MatchLabel{
		Size: int(size),
		Open: 1,
	}
	if fast {
		Label.Fast = 1
	}

	labelJSON, err := json.Marshal(Label)
	if err != nil {
		logger.WithField("error", err).Error("match init failed")
		labelJSON = []byte("{}")
	}
	// logger.Debug("*****创建房间", labelJSON)
	return &model.MatchState{
		Random:                 rand.New(rand.NewSource(time.Now().UnixNano())),
		Label:                  Label,
		NextGameRemainingTicks: nextStartSec * tickRate,

		Presences: make(map[string]runtime.Presence, 4),
	}, tickRate, string(labelJSON)
}

func (m *MatchHandler) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	s := state.(*model.MatchState)

	// Check if it's a user attempting to rejoin after a disconnect.
	if presence, ok := s.Presences[presence.GetUserId()]; ok {
		if presence == nil {
			// User rejoining after a disconnect.
			s.JoinsInProgress++
			return s, true, ""
		} else {
			// User attempting to join from 2 different devices at the same time.
			return s, false, "already joined"
		}
	}

	// Check if match is full.
	if len(s.Presences)+s.JoinsInProgress >= s.Label.Size {
		return s, false, "match full"
	}
	// New player attempting to connect.
	s.JoinsInProgress++
	return s, true, ""
}

func (m *MatchHandler) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, Presences []runtime.Presence) interface{} {

	s := state.(*model.MatchState)

	for _, presence := range Presences {
		s.EmptyTicks = 0
		s.Presences[presence.GetUserId()] = presence
		s.JoinsInProgress--

		// Check if we must send a message to this user to update them on the current game state.
		var opCode api.OpCode
		var msg proto.Message
		if s.Playing {
			// There's a game still currently in progress, the player is re-joining after a disconnect. Give them a state update.
			opCode = api.OpCode_OPCODE_UPDATE
			msg = &api.Update{
				PlayerInfos:     s.PlayerInfos,
				CurrentPlayerId: s.CurrentPlayerID,
				Message:         "reconnect",
			}
		}

		// Send a message to the user that just joined, if one is needed based on the logic above.
		if msg != nil {
			buf, err := global.Marshaler.Marshal(msg)
			if err != nil {
				logger.Error("error encoding message: %v", err)
			} else {
				_ = dispatcher.BroadcastMessage(int64(opCode), buf, []runtime.Presence{presence}, nil, true)
			}
		}
	}

	// Check if match was open to new players, but should now be closed.
	if len(s.Presences) >= s.Label.Size && s.Label.Open != 0 {
		s.Label.Open = 0
		if labelJSON, err := json.Marshal(s.Label); err != nil {
			logger.Error("error encoding Label: %v", err)
		} else {
			if err := dispatcher.MatchLabelUpdate(string(labelJSON)); err != nil {
				logger.Error("error updating Label: %v", err)
			}
		}
	}
	return s
}

func (m *MatchHandler) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, Presences []runtime.Presence) interface{} {
	s := state.(*model.MatchState)
	for _, presence := range Presences {
		logger.Warn("%v leaves match", presence.GetUserId())
		s.Presences[presence.GetUserId()] = nil
	}

	return s
}

func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {

	s := state.(*model.MatchState)
	//关闭房间
	if len(s.Presences)+s.JoinsInProgress == 0 {
		s.EmptyTicks++
		if s.EmptyTicks >= maxEmptySec*tickRate {
			// Match has been empty for too long, close it.
			logger.Info("closing idle match")
			return nil
		}
	}

	// If there's no game in progress check if we can (and should) start one!
	if !s.Playing {
		// Between games any disconnected users are purged, there's no in-progress game for them to return to anyway.
		for userID, presence := range s.Presences {
			if presence == nil {
				delete(s.Presences, userID)
			}
		}

		// Check if we need to update the Label so the match now advertises itself as open to join.
		if len(s.Presences) < s.Label.Size && s.Label.Open != 1 {
			s.Label.Open = 1
			if labelJSON, err := json.Marshal(s.Label); err != nil {
				logger.Error("error encoding Label: %v", err)
			} else {
				if err := dispatcher.MatchLabelUpdate(string(labelJSON)); err != nil {
					logger.Error("error updating Label: %v", err)
				}
			}
		}

		// Check if we have enough players to start a game.
		if len(s.Presences) < s.Label.Size {
			return s
		}
		// Check if enough time has passed since the last game.
		if s.NextGameRemainingTicks > 0 {
			if s.NextGameRemainingTicks%tickRate == 0 {
				buf, err := global.Marshaler.Marshal(&api.ReadyToStart{
					NextGameStart: s.NextGameRemainingTicks / tickRate,
				})
				if err != nil {
					logger.Error("error encoding message: %v", err)
				} else {
					_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_READY_START), buf, nil, nil, true)
				}
			}
			s.NextGameRemainingTicks--
			return s
		}
		// We can start a game!
		s.Playing = true
		//初始化游戏
		m.service.StartMatch(dispatcher, logger, s, tickRate, turnTimeNormalSec)
		return s
	}

	// There's a game in progress. Check for input, update match state, and send messages to clients.
	for _, message := range messages {
		serv.Dispatch(message, s, dispatcher)
	}
	// Keep track of the time remaining for the player to submit their move. Idle players forfeit.
	if s.Playing {
		s.DeadlineRemainingTicks--
		if s.DeadlineRemainingTicks%tickRate == 0 {
			buf, _ := global.Marshaler.Marshal(&api.Tick{Deadline: s.DeadlineRemainingTicks / tickRate})
			_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_TICK), buf, nil, nil, true)
		}
		if s.DeadlineRemainingTicks <= 0 {
			// The player has run out of time to submit their move.
			serv.DefaultAction(dispatcher, s)
		}
	}
	return s
}

func (m *MatchHandler) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	return state, ""
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	return state
}
