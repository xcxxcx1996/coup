package service

import (
	"fmt"
	"log"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

type DrawCoin struct {
	message runtime.MatchData
	coins   int32
	IsDeny  bool
}

func (a DrawCoin) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {

	msg := &api.GetCoin{}
	myTurn := message.GetUserId() == state.CurrentPlayerID

	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	if err != nil || !myTurn {
		// Client sent bad data.
		log.Printf("错误的参数:%v , 不是我的回合:%v", err, myTurn)
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	a.coins = msg.Coins
	a.message = message
	state.Actions.Push(a)
	var buf []byte
	// question状态
	if a.coins == 2 {
		buf, _ = global.Marshaler.Marshal(&api.Info{Info: fmt.Sprintf("%v想拿2块钱", message.GetUsername())})
		state.EnterDenyMoney()
	} else {
		buf, _ = global.Marshaler.Marshal(&api.Info{Info: fmt.Sprintf("%v想拿1块钱", message.GetUsername())})
	}
	log.Println("拿钱")
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
}

// 正常拿钱没有人会质疑，只有公爵阻止,所以不用写
func (a DrawCoin) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {

}

// 阻止失败，开始拿钱
func (a DrawCoin) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	state.PlayerInfos[state.CurrentPlayerID].Coins += int64(a.coins)
	info := fmt.Sprintf("%v成功拿了钱", a.message.GetUsername())
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
}

// 行动被阻止
func (a DrawCoin) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	state.Actions.Pop()
	state.NextTurn()
}

func (a DrawCoin) GetRole() api.Role {
	return api.Role_UNROLE
}
