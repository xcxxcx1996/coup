package main

import (
	"fmt"

	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
)

func main() {
	global.Init()
	t := &api.DenyMoney{}
	global.Unmarshaler.Unmarshal(bytes, t)
	fmt.Printf("t:%v", t)
}
