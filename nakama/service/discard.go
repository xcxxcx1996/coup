package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

//玩家质疑成功或者失败弃牌的真实函数,传入那张牌
func (serv *MatchService) Discard(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	//
	msg := &api.Discard{}
	valid := ValidAction(state, message, api.State_DISCARD, msg)
	// 推进行动
	if !valid {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	var discard *api.Card
	cards := state.PlayerInfos[message.GetUserId()].Cards
	for i, c := range cards {
		if c.Id == msg.CardId {
			// discard = c
			cards = append(cards[:i], cards[i+1:]...)
		}
	}
	info := fmt.Sprintf("%v discard the %v", message.GetUsername(), discard.Role)
	SendNotification(info, dispatcher)

	// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: info})
	// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	// 判断玩家是否会死
	alive, _ := state.Alive(message.GetUserId())
	if !alive {
		state.EliminatePlayer(message.GetUserId())
		info := fmt.Sprintf("%v was eliminated", message.GetUserId())
		SendNotification(info, dispatcher)
		// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: info})
		// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
		//
		buf, _ := global.Marshaler.Marshal(&api.Dead{Player: state.PlayerInfos[message.GetUserId()]})
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DEAD), buf, nil, nil, true)
		// 判断冠军
		if len(state.PlayerSequence) == 1 {
			buf, _ = global.Marshaler.Marshal(&api.Done{Winner: state.PlayerInfos[state.PlayerSequence[0]]})
			_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DONE), buf, nil, nil, true)
		}
	}
	if state.ActionComplete {
		state.NextTurn()
	}
}
