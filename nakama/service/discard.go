package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

// func (serv *MatchService) Discarding(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {

// 	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)

// }
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
	cards := state.PlayerInfos[msg.PlayerId].Cards
	for i, c := range cards {
		if c.Id == msg.CardId {
			discard = c
			cards = append(cards[:i], cards[i+1:]...)
		}
	}
	info := fmt.Sprintf("%v弃牌 %v", message.GetUsername(), discard.Role)
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
}
