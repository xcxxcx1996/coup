package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

// func (serv *MatchService) Discarding(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {

// 	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)

// }
//玩家质疑成功或者失败弃牌的真实函数,传入那张牌
func (serv *MatchService) Discard(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	//
	msg := &api.Discard{}
	if err := serv.Unmarshaler.Unmarshal(message.GetData(), msg); err != nil {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	cards := state.PlayerInfos[msg.PlayerId].Cards
	for i, c := range cards {
		if c.Id == msg.CardId {
			cards = append(cards[:i], cards[i+1:]...)
		}
	}
	//如果行为完成了，那么就是下一个人，没有完成，那么就是

}

// 玩家进入弃牌阶段
func (serv *MatchService) EnterDicardState(playerID string, state *model.MatchState) {
	//下一个用户进入question状态
	for _, p := range state.PlayerInfos {
		if playerID == p.Id {
			p.State = api.State_DISCARD
		} else {
			p.State = api.State_IDLE
		}
	}
}
