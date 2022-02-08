package main

import (
	"log"

	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func main() {
	global.Init()
	buf, _ := global.Marshaler.Marshal(&api.Message{Info: "aaa"})
	valid(buf, nil)
	msg := &api.Message{}
	// 数据是否准确
	log.Println(err)
}

func valid(buf []byte, m protoreflect.ProtoMessage) (ok bool) {
	err := global.Unmarshaler.Unmarshal(buf, m)
	if err != nil {
		return false
	}
	return true

}
