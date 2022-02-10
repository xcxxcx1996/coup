package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

func (serv *MatchService) Questioning(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {

	msg := &api.Question{}

	if err = ValidAction(state, message, api.State_QUESTION, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// 质疑
	if msg.IsQuestion {
		info := fmt.Sprintf(`<p><span style="color:red;">%v</span> questioned the action.</p >`, message.GetUsername())
		SendNotification(info, dispatcher)
		if state.ValidQuestion() {
			// 质疑成功让某人弃牌
			action, _ := state.Actions.Pop()
			info := fmt.Sprintf(`<p><span style="color:red;">%v</span> question <span style="color:green;">successful</span> and <span style="color:red;">%v</span> is discarding.</p >`, message.GetUsername(), state.GetPlayerNameByID(action.GetActor()))
			SendNotification(info, dispatcher)
			action.Stop(dispatcher, state)
			state.EnterDicardState(action.GetActor())

		} else {
			// 质疑失败 自己进入弃牌
			info := fmt.Sprintf(`<p><span style="color:red;">%v</span> question <span style="color:red;">failed</span> and <span style="color:red;">%v</span> is discarding</p >`, message.GetUsername(), message.GetUsername())
			SendNotification(info, dispatcher)
			state.EnterDicardState(message.GetUserId())
		}
		return
	} else {
		// 不质疑
		// 如果下一个是当前行动人，则说明循环了一圈，退出Question，进入刺杀阶段"

		info := fmt.Sprintf(`<p><span style="color:red;">%v</span> didn't question.</p>`, message.GetUsername())
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
