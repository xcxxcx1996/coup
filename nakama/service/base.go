package service

import (
	"log"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MatchService struct{}

func New() (s *MatchService) {
	return &MatchService{}
}

func (serv *MatchService) Dispatch(message runtime.MatchData, state *model.MatchState, dispatcher runtime.MatchDispatcher) {
	switch message.GetOpCode() {
	// 刺客刺杀
	case int64(api.OpCode_OPCODE_ASSASSIN):
		action := Assassin{}
		action.Start(dispatcher, message, state)
	// 大使换牌
	case int64(api.OpCode_OPCODE_CHANGE_CARD):
		action := ChangeCard{}
		action.Start(dispatcher, message, state)
	// 大使选牌
	case int64(api.OpCode_OPCODE_CHOOSE_CARD):
		serv.CompleteChangeCard(dispatcher, message, state)
	// 政变
	case int64(api.OpCode_OPCODE_COUP):
		serv.Coup(dispatcher, message, state)
	// 女王阻止刺杀
	case int64(api.OpCode_OPCODE_DENY_KILL):
		action := DenyAssassian{}
		action.Start(dispatcher, message, state)
	// 男爵阻止拿牌
	case int64(api.OpCode_OPCODE_DENY_MONEY):
		action := DenyMoney{}
		action.Start(dispatcher, message, state)
	// 男爵拿钱
	case int64(api.OpCode_OPCODE_DRAW_THREE_COINS):
		action := DrawThreeCoins{}
		action.Start(dispatcher, message, state)
	//偷钱
	case int64(api.OpCode_OPCODE_STEAL_COINS):
		action := Steal{}
		action.Start(dispatcher, message, state)
	//阻止偷钱
	case int64(api.OpCode_OPCODE_DENY_STEAL):
		action := DenySteal{}
		action.Start(dispatcher, message, state)
	// 普通玩家拿钱
	case int64(api.OpCode_OPCODE_DRAW_COINS):
		action := DrawCoin{}
		action.Start(dispatcher, message, state)
	// 普通玩家弃牌
	case int64(api.OpCode_OPCODE_DISCARD):
		serv.Discard(dispatcher, message, state)
	// 普通玩家质疑
	case int64(api.OpCode_OPCODE_QUESTION):
		serv.Questioning(dispatcher, message, state)
	default:
	}
	var opCode api.OpCode
	var outgoingMsg proto.Message
	opCode = api.OpCode_OPCODE_UPDATE
	outgoingMsg = &api.Update{
		PlayerInfos:     state.PlayerInfos,
		CurrentPlayerId: state.CurrentPlayerID,
		Message:         "update",
	}

	buf, err := global.Marshaler.Marshal(outgoingMsg)
	if err != nil {
		log.Printf("error encoding message: %v", err)
	} else {
		_ = dispatcher.BroadcastMessage(int64(opCode), buf, nil, nil, true)
	}
	state.ResetDeadLine()
}

// 默认行为
func (serv *MatchService) DefaultAction(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	log.Printf("state:   %v", state.State)
	var message DefaultActionData
	switch state.State {
	case api.State_QUESTION:
		data, _ := global.Marshaler.Marshal(&api.Question{IsQuestion: false})
		message = DefaultActionData{OpCode: int64(api.OpCode_OPCODE_QUESTION), Data: data, Presence: state.Presences[state.Currentquestioner]}
	case api.State_DISCARD:
		data, _ := global.Marshaler.Marshal(&api.Discard{CardId: state.PlayerInfos[state.CurrentDiscarder].Cards[0].Id})
		message = DefaultActionData{OpCode: int64(api.OpCode_OPCODE_DRAW_COINS), Data: data, Presence: state.Presences[state.CurrentDiscarder]}
	case api.State_CHOOSE_CARD:
		cards := []string{}
		for _, v := range state.PlayerInfos[state.CurrentDiscarder].Cards {
			cards = append(cards, v.Id)
		}
		data, _ := global.Marshaler.Marshal(&api.ChangeCardResponse{Cards: cards})
		message = DefaultActionData{OpCode: int64(api.OpCode_OPCODE_CHOOSE_CARD), Data: data, Presence: state.Presences[state.CurrentDiscarder]}
	case api.State_START:
		data, _ := global.Marshaler.Marshal(&api.GetCoin{Coins: 1})
		message = DefaultActionData{OpCode: int64(api.OpCode_OPCODE_DRAW_COINS), Data: data, Presence: state.Presences[state.CurrentPlayerID]}
	case api.State_DENY_MONEY:
		data, _ := global.Marshaler.Marshal(&api.Deny{IsDeny: false})
		message = DefaultActionData{OpCode: int64(api.OpCode_OPCODE_DENY_MONEY), Data: data, Presence: state.Presences[state.CurrentDenyer]}
	case api.State_DENY_STEAL:
		data, _ := global.Marshaler.Marshal(&api.Deny{IsDeny: false})
		message = DefaultActionData{OpCode: int64(api.OpCode_OPCODE_DENY_STEAL), Data: data, Presence: state.Presences[state.CurrentDenyer]}
	case api.State_DENY_ASSASSIN:
		data, _ := global.Marshaler.Marshal(&api.Deny{IsDeny: false})
		message = DefaultActionData{OpCode: int64(api.OpCode_OPCODE_DENY_KILL), Data: data, Presence: state.Presences[state.CurrentDenyer]}

	}
	serv.Dispatch(message, state, dispatcher)
}

func (serv *MatchService) StartMatch(dispatcher runtime.MatchDispatcher, logger runtime.Logger, s *model.MatchState, tickRate int64, turnTime int64) *model.MatchState {
	// Notify the players a new game has started.
	s.Init(turnTime * tickRate)
	s.State = api.State_START
	buf, err := global.Marshaler.Marshal(&api.Start{
		PlayerInfos:     s.PlayerInfos,
		CurrentPlayerId: s.CurrentPlayerID,
		Message:         s.Message,
		Deadline:        s.DeadlineRemainingTicks / tickRate,
	})
	if err != nil {
		logger.Error("error encoding message: %v", err)
	} else {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_START), buf, nil, nil, true)
	}
	return s
}

func SendNotification(msg string, dispatcher runtime.MatchDispatcher) {
	buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: msg})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
}

func ValidAction(state *model.MatchState, message runtime.MatchData, allowState api.State, msg protoreflect.ProtoMessage) (ok bool) {
	if msg != nil {
		err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
		if err != nil {
			log.Println("wrong match data")
			// Client sent bad data.
			return false
		}
	}
	if allowState != state.State {
		return false
	}
	switch allowState {
	case api.State_START:
		ok = message.GetUserId() == state.CurrentPlayerID
	case api.State_CHOOSE_CARD:
		ok = message.GetUserId() == state.CurrentPlayerID
	case api.State_DISCARD:
		ok = message.GetUserId() == state.CurrentDiscarder
	case api.State_QUESTION:
		ok = message.GetUserId() == state.Currentquestioner
	default:
		ok = message.GetUserId() == state.CurrentDenyer
	}
	return
}
