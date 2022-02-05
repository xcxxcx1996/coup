package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) Questioning(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	msg := &api.Question{}
	err := serv.Unmarshaler.Unmarshal(message.GetData(), msg)

	if err != nil {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
	}
	// 质疑的判断
	if msg.IsQuestion{

		if judgeQuestion() {
			// 质疑成功
			// 某些效果
		} else {
			//质疑失败 丢弃某张卡牌
			// 自己进入弃牌
		}
		// 质疑阶段结束

		return
	} else {

	}

	//

	// 不质疑

	// 如果下一个质疑玩家是当前玩家则退出

	// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)
}

func (serv *MatchService) AfterQuestion(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	//如果

	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)

}
func judgeQuestion() bool {
	return true
}

func (serv *MatchService) EnterQuestionState(playerID string, state *model.MatchState) {
	//下一个用户进入question状态
	for _, p := range state.PlayerInfos {
		if playerID == p.Id {
			p.State = api.State_QUESTION
		} else {
			p.State = api.State_IDLE
		}
	}
}
