package service

import (
	"fmt"
	"log"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

type Assassin struct {
	Assassinor   string
	Assassinated string
}

func (a Assassin) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	// 获得信息、核验

	msg := &api.Assassin{}
	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	myTurn := message.GetUserId() == state.CurrentPlayerID
	if err != nil || !myTurn {
		// Client sent bad data.
		log.Printf("错误的参数:%v , 不是我的回合:%v", err, myTurn)
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// 推进行动
	a.Assassinated = msg.PlayerId
	a.Assassinor = message.GetUserId()
	state.Actions.Push(a)

	//后处理

	info := fmt.Sprintf("%v want to assassin %v ", message.GetUsername(), state.GetPlayerNameByID(a.Assassinated))
	SendNotification(info, dispatcher)
	state.EnterQuestion()
}

// 1. 没人质疑，2. 有人质疑，但失败
func (a Assassin) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("question end，deny start")
	SendNotification(info, dispatcher)

	state.EnterDenyAssassin(a.Assassinated)
}

// 玩家进入刺杀阶段
func (a Assassin) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("deny end, action execute")
	SendNotification(info, dispatcher)

	state.EnterDicardState(a.Assassinated)
}

// 行动被停止，删除，然后跳过
func (a Assassin) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("assassin was stopped")
	SendNotification(info, dispatcher)

	state.Actions.Pop()
	state.NextTurn()
}

func (a Assassin) GetRole() api.Role {
	return api.Role_ASSASSIN
}
