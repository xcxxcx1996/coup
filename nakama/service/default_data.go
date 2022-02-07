package service

import "github.com/heroiclabs/nakama-common/runtime"

type DefaultActionData struct {
	runtime.Presence
	OpCode int64
	Data   []byte
}

func (d DefaultActionData) GetOpCode() int64 {
	return d.OpCode
}

func (d DefaultActionData) GetData() []byte {
	return d.Data
}

func (d DefaultActionData) GetReliable() bool {
	return false
}

func (d DefaultActionData) GetReceiveTime() int64 {
	return int64(0)
}
