package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
	"github.com/xcxcx1996/coup/service"
)

const (
	moduleName = "coup"

	tickRate = 5

	maxEmptySec = 30

	// delayBetweenGamesSec = 5
	turnTimeFastSec   = 10
	turnTimeNormalSec = 20
)

// Compile-time check to make sure all required functions are implemented.
var serv = service.New()
var _ runtime.Match = &MatchHandler{
	service: serv,
}

type MatchHandler struct {
	marshaler   *protojson.MarshalOptions
	unmarshaler *protojson.UnmarshalOptions
	service     *service.MatchService
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	fast, ok := params["fast"].(bool)
	if !ok {
		logger.Error("invalid match init parameter \"fast\"")
		return nil, 0, ""
	}

	Label := &model.MatchLabel{
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
	logger.Debug("*****创建房间", labelJSON)
	return &model.MatchState{
		Random: rand.New(rand.NewSource(time.Now().UnixNano())),
		Label:  Label,

		Presences: make(map[string]runtime.Presence, 4),
	}, tickRate, string(labelJSON)
}

func (m *MatchHandler) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	s := state.(*model.MatchState)

	//

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
	if len(s.Presences)+s.JoinsInProgress >= 3 {
		return s, false, "match full"
	}
	// New player attempting to connect.
	s.JoinsInProgress++
	return s, true, ""
}

func (m *MatchHandler) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, Presences []runtime.Presence) interface{} {

	s := state.(*model.MatchState)

	// t := time.Now().UTC()

	for _, presence := range Presences {
		s.EmptyTicks = 0
		s.Presences[presence.GetUserId()] = presence
		s.JoinsInProgress--

		// Check if we must send a message to this user to update them on the current game state.
		var opCode api.OpCode
		var msg proto.Message

		// There's no game in progress but we still have a completed game that the user was part of.
		// They likely disconnected before the game ended, and have since forfeited because they took too long to return.
		//todo aaa
		// 	else if s.board != nil && s.marks != nil && s.marks[presence.GetUserId()] > api.Mark_MARK_UNSPECIFIED {
		// 	opCode = api.OpCode_OPCODE_DONE
		// 	msg = &api.Done{
		// 		Board:           s.board,
		// 		Winner:          s.winner,
		// 		WinnerPositions: s.winnerPositions,
		// 		NextGameStart:   t.Add(time.Duration(s.nextGameRemainingTicks/tickRate) * time.Second).Unix(),
		// 	}
		// }

		// Send a message to the user that just joined, if one is needed based on the logic above.
		if msg != nil {
			buf, err := m.marshaler.Marshal(msg)
			if err != nil {
				logger.Error("error encoding message: %v", err)
			} else {
				_ = dispatcher.BroadcastMessage(int64(opCode), buf, []runtime.Presence{presence}, nil, true)
			}
		}
	}

	// Check if match was open to new players, but should now be closed.
	if len(s.Presences) >= 3 && s.Label.Open != 0 {
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

	t := time.Now().UTC()

	// If there's no game in progress check if we can (and should) start one!
	if !s.Playing {
		// Between games any disconnected users are purged, there's no in-progress game for them to return to anyway.
		for userID, presence := range s.Presences {
			if presence == nil {
				delete(s.Presences, userID)
			}
		}

		// Check if we need to update the Label so the match now advertises itself as open to join.
		if len(s.Presences) < 2 && s.Label.Open != 1 {
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
		if len(s.Presences) < 2 {
			return s
		}

		// Check if enough time has passed since the last game.
		// if s.nextGameRemainingTicks > 0 {
		// 	s.nextGameRemainingTicks--
		// 	return s
		// }

		// We can start a game! Set up the game state and assign the marks to each player.
		s.Playing = true
		//初始化游戏
		//todo

		// Notify the players a new game has started.
		buf, err := m.marshaler.Marshal(&api.Start{
			// Board:    s.board,
			// Marks:    s.marks,
			// Mark:     s.mark,
			Deadline: t.Add(time.Duration(s.DeadlineRemainingTicks/tickRate) * time.Second).Unix(),
		})
		if err != nil {
			logger.Error("error encoding message: %v", err)
		} else {
			_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_START), buf, nil, nil, true)
		}
		return s
	}

	// There's a game in progress. Check for input, update match state, and send messages to clients.
	for _, message := range messages {
		switch api.OpCode(message.GetOpCode()) {
		case api.OpCode_OPCODE_ASSASSIN:
			m.service.BeforeAssassin(ctx, dispatcher, s, message)
			log.Println("this is Assassin")
		case api.OpCode_OPCODE_DENY_KILL:

			log.Println("this is Assassin")
		case api.OpCode_OPCODE_DISCARD:
			log.Println("this is Assassin")
		case api.OpCode_OPCODE_DENY_STEAL:
			log.Println("this is Assassin")
		case api.OpCode_OPCODE_COUP:
			log.Println("this is Assassin")
		case api.OpCode_OPCODE_DRAW_CARD:
			log.Println("this is Assassin")
		case api.OpCode_OPCODE_DRAW_COINS:
			log.Println("this is Assassin")
		case api.OpCode_OPCODE_CHANGE_CARD:
			log.Println("this is Assassin")
		case api.OpCode_OPCODE_DRAW_THREE_COINS:
			log.Println("this is Assassin")
		case api.OpCode_OPCODE_STEAL_CARD:
			log.Println("this is Assassin")

		case api.OpCode_OPCODE_QUESTION:
			log.Println("this is Assassin")
		case api.OpCode_OPCODE_DONE:
			log.Println("this is Assassin")

		// case api.OpCode_OPCODE_MOVE:
		// 	mark := s.marks[message.GetUserId()]
		// 	if s.mark != mark {
		// 		// It is not this player's turn.
		// 		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, []runtime.Presence{message}, nil, true)
		// 		continue
		// 	}

		// 	msg := &api.Move{}
		// 	err := m.unmarshaler.Unmarshal(message.GetData(), msg)
		// 	if err != nil {
		// 		// Client sent bad data.
		// 		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, []runtime.Presence{message}, nil, true)
		// 		continue
		// 	}
		// 	if msg.Position < 0 || msg.Position > 8 || s.board[msg.Position] != api.Mark_MARK_UNSPECIFIED {
		// 		// Client sent a position outside the board, or one that has already been played.
		// 		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, []runtime.Presence{message}, nil, true)
		// 		continue
		// 	}

		// 	// Update the game state.
		// 	s.board[msg.Position] = mark
		// 	switch mark {
		// 	case api.Mark_MARK_X:
		// 		s.mark = api.Mark_MARK_O
		// 	case api.Mark_MARK_O:
		// 		s.mark = api.Mark_MARK_X
		// 	}
		// 	s.DeadlineRemainingTicks = calculateDeadlineTicks(s.Label)

		// 	// Check if game is over through a winning move.
		// //winCheck:
		// 	for _, winningPosition := range winningPositions {
		// 		for _, position := range winningPosition {
		// 			if s.board[position] != mark {
		// 				continue winCheck
		// 			}
		// 		}

		// 		// Update state to reflect the winner, and schedule the next game.
		// 		s.winner = mark
		// 		s.winnerPositions = winningPosition
		// 		s.Playing = false
		// 		s.DeadlineRemainingTicks = 0
		// 		s.nextGameRemainingTicks = delayBetweenGamesSec * tickRate
		// 	}
		// 	// Check if game is over because no more moves are possible.
		// 	tie := true
		// 	for _, mark := range s.board {
		// 		if mark == api.Mark_MARK_UNSPECIFIED {
		// 			tie = false
		// 			break
		// 		}
		// 	}
		// 	if tie {
		// 		// Update state to reflect the tie, and schedule the next game.
		// 		s.Playing = false
		// 		s.DeadlineRemainingTicks = 0
		// 		s.nextGameRemainingTicks = delayBetweenGamesSec * tickRate
		// 	}

		// 	var opCode api.OpCode
		// 	var outgoingMsg proto.Message
		// 	if s.Playing {
		// 		opCode = api.OpCode_OPCODE_UPDATE
		// 		outgoingMsg = &api.Update{
		// 			Board:    s.board,
		// 			Mark:     s.mark,
		// 			Deadline: t.Add(time.Duration(s.DeadlineRemainingTicks/tickRate) * time.Second).Unix(),
		// 		}
		// 	} else {
		// 		opCode = api.OpCode_OPCODE_DONE
		// 		outgoingMsg = &api.Done{
		// 			Board:           s.board,
		// 			Winner:          s.winner,
		// 			WinnerPositions: s.winnerPositions,
		// 			NextGameStart:   t.Add(time.Duration(s.nextGameRemainingTicks/tickRate) * time.Second).Unix(),
		// 		}
		// 	}

		// 	buf, err := m.marshaler.Marshal(outgoingMsg)
		// 	if err != nil {
		// 		logger.Error("error encoding message: %v", err)
		// 	} else {
		// 		_ = dispatcher.BroadcastMessage(int64(opCode), buf, nil, nil, true)
		// 	}
		default:
			// No other opcodes are expected from the client, so automatically treat it as an error.
			_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DONE), nil, []runtime.Presence{message}, nil, true)
		}
	}

	// Keep track of the time remaining for the player to submit their move. Idle players forfeit.
	if s.Playing {
		s.DeadlineRemainingTicks--
		if s.DeadlineRemainingTicks <= 0 {
			// The player has run out of time to submit their move.
			log.Println("DeadlineRemainingTicks")
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

// func calculateDeadlineTicks(l *MatchLabel) int64 {
// 	if l.Fast == 1 {
// 		return turnTimeFastSec * tickRate
// 	} else {
// 		return turnTimeNormalSec * tickRate
// 	}
// }
