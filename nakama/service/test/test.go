package main

import (
	"log"
	"unsafe"

	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
)

func BytesToStringFast(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func main() {
	// tickRate := int64(10)
	global.Init()
	// t := time.Now().UTC()
	buf, _ := global.Marshaler.Marshal(&api.Question{
		IsQuestion: false,
	})
	log.Printf("buf:%v", buf)
	log.Printf("string:%v", BytesToStringFast(buf))
	var a = api.Question{}
	global.Unmarshaler.Unmarshal(buf, &a)
	log.Printf("struct:%v", a)
}
