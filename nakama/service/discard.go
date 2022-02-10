package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	model "github.com/xcxcx1996/coup/state"
)

//玩家质疑成功或者失败弃牌的真实函数,传入那张牌
func (serv *MatchService) Discard(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	// 验证
	msg := &api.Discard{}
	if err = ValidAction(state, message, api.State_DISCARD, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	discard, err := state.DeleteCard(msg.CardId, message.GetUserId())
	if err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	info := fmt.Sprintf(`<p><span style={{ color: "red" }}>%v</span> discard the <span style={{ color: "red" }}>%v</span>.</p >`, message.GetUsername(), discard.Role)
	SendNotification(info, dispatcher)

	// 判断玩家是否会死
	alive, _ := state.Alive(message.GetUserId())
	if !alive {
		state.EliminatePlayer(message.GetUserId())
		info := fmt.Sprintf(`<p><span style={{ color: "red" }}>%v</span> was dead.</p >`, message.GetUsername())
		SendNotification(info, dispatcher)
		buf, _ := global.Marshaler.Marshal(&api.Dead{Player: state.PlayerInfos[message.GetUserId()]})
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DEAD), buf, nil, nil, true)
		// 判断冠军
		if len(state.PlayerSequence) == 1 {
			state.Playing = false
			state.ResetNextMartch()
			info := fmt.Sprintln("<p>Match end!!!</p>")
			SendNotification(info, dispatcher)
			buf, _ = global.Marshaler.Marshal(&api.Done{Winner: state.PlayerInfos[state.PlayerSequence[0]]})
			_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DONE), buf, nil, nil, true)
		}
		state.NextTurn()
		// 角色死亡quit
		return
	}
	action, err := state.Actions.Pop()

	if err != nil {
		state.NextTurn()
		return nil
	}
	ass, ok := action.(Assassin)
	if ok {
		err = ass.AfterDeny(dispatcher, state)
		return
	}
	if !state.ActionComplete {
		action.AfterQuestion(dispatcher, state)
		return
	}
	return
}
