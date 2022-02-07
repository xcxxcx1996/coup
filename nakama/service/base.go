package service

import (
	"log"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
	"google.golang.org/protobuf/proto"
)

type MatchService struct{}

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
	// 普通玩家拿牌
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
}

// 默认行为
func (serv *MatchService) DefaultAction(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
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
	default:
		data, _ := global.Marshaler.Marshal(&api.Deny{IsDeny: false})
		message = DefaultActionData{OpCode: int64(api.OpCode_OPCODE_DENY_KILL), Data: data, Presence: state.Presences[state.CurrentDenyer]}
	}
	serv.Dispatch(message, state, dispatcher)
}

func New() (s *MatchService) {
	return &MatchService{}
}
