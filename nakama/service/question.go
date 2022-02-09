package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

func (serv *MatchService) Questioning(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	// 质疑

	msg := &api.Question{}

	if err = ValidAction(state, message, api.State_QUESTION, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// 质疑的判断
	if msg.IsQuestion {
		info := fmt.Sprintf("%v questioned the action.", message.GetUsername())
		SendNotification(info, dispatcher)
		if state.ValidQuestion() {
			// 质疑成功让某人弃牌
			info := fmt.Sprintf("%v question successful and %v is discarding.", message.GetUsername(), state.GetPlayerNameByID(state.CurrentPlayerID))
			SendNotification(info, dispatcher)
			action, _ := state.Actions.Pop()
			// 中止
			action.Stop(dispatcher, state)
			state.EnterDicardState(action.GetActor())

		} else {
			// 质疑失败 自己进入弃牌
			info := fmt.Sprintf("%v question failed, %v is discarding.", message.GetUsername(), message.GetUsername())
			SendNotification(info, dispatcher)
			state.EnterDicardState(message.GetUserId())
		}
		return
	} else {
		// 不质疑
		// 如果下一个是当前行动人，则说明循环了一圈，退出Question，进入刺杀阶段
		info := fmt.Sprintf("%v didn't question.", message.GetUsername())
		SendNotification(info, dispatcher)
		end := state.NextQuestionor()
		if end {
			action, _ := state.Actions.Pop()
			err = action.AfterQuestion(dispatcher, state)
			return err
		}
	}
	return
}
