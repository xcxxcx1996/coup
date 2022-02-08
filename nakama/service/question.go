package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) Questioning(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	// 质疑

	msg := &api.Question{}

	myTurn := message.GetUserId() == state.CurrentPlayerID

	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	if err != nil || !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// 质疑的判断
	if msg.IsQuestion {
		card_id, valid := state.ValidQuestion()
		info := fmt.Sprintf("%v questioned the action", message.GetUsername())
		SendNotification(info, dispatcher)
		// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: fmt.Sprintf("%v 质询了该行动", message.GetUsername())})
		// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
		if valid {
			// 质疑成功让某人弃牌
			info := fmt.Sprintf("%v question successful", message.GetUsername())
			SendNotification(info, dispatcher)
			// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: fmt.Sprintf("%v 质询成功", message.GetUsername())})
			// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
			state.EnterDicardState(state.CurrentPlayerID)
			action, _ := state.Actions.Last()
			// 中止
			action.Stop(dispatcher, state)
		} else {
			//质疑失败 自己进入弃牌
			info := fmt.Sprintf("%v question failed", message.GetUsername())
			SendNotification(info, dispatcher)
			// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: fmt.Sprintf("%v 质询失败", message.GetUsername())})
			// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
			state.EnterDicardState(message.GetUserId())
			serv.ChangeSingleCard(card_id, state.CurrentPlayerID, state)
			//继续他人的行动
			action, _ := state.Actions.Last()
			action.AfterQuestion(dispatcher, state)
			// 然后换一张牌
		}
		return
	} else {
		// 不质疑
		// 如果下一个是当前行动人，则说明循环了一圈，退出Question，进入刺杀阶段
		info := fmt.Sprintf("%v didn't question", message.GetUsername())
		SendNotification(info, dispatcher)
		// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: fmt.Sprintf("%v 不质询", message.GetUsername())})
		// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
		end := state.NextQuestionor()
		if end {
			action, _ := state.Actions.Last()
			action.AfterQuestion(dispatcher, state)
			return
		}

	}
}
