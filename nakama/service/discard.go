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
	myTurn := message.GetUserId() == state.CurrentPlayerID
	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)

	if err != nil || !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	var discard *api.Card
	cards := state.PlayerInfos[message.GetUserId()].Cards
	for i, c := range cards {
		if c.Id == msg.CardId {
			discard = c
			cards = append(cards[:i], cards[i+1:]...)
		}
	}
	info := fmt.Sprintf("%v弃牌 %v", message.GetUsername(), discard.Role)
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	// 判断玩家是否会死
	alive, _ := state.Alive(message.GetUserId())
	if !alive {
		state.EliminatePlayer(message.GetUserId())
		info := fmt.Sprintf("%v被淘汰了", message.GetUserId())
		buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
		//
		buf, _ = global.Marshaler.Marshal(&api.Dead{Player: state.PlayerInfos[message.GetUserId()]})
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DEAD), buf, nil, nil, true)
		// 判断冠军
		if len(state.PlayerSequence) == 1 {
			buf, _ = global.Marshaler.Marshal(&api.Done{Winner: state.PlayerInfos[state.PlayerSequence[0]]})
			_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DONE), buf, nil, nil, true)
		}
	}

}
