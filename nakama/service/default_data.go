package service

import (
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
)

type DefaultPresence struct {
	runtime.PresenceMeta
	PlayerInfo *api.PlayerInfo
}

func (d DefaultPresence) GetUserId() string {
	return d.PlayerInfo.Id
}
func (d DefaultPresence) GetUsername() string {
	return d.PlayerInfo.Name
}

func (d DefaultPresence) GetSessionId() string {
	return ""
}
func (d DefaultPresence) GetNodeId() string {
	return ""
}

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
