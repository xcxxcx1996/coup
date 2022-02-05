package main

import (
	"log"
	"time"

	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
	"github.com/xcxcx1996/coup/service"
)

func main() {
	tickRate := int64(10)
	s := &model.MatchState{
		NextGameRemainingTicks: 40,
	}
	service := service.New()
	t := time.Now().UTC()
	buf, err := service.Marshaler.Marshal(&api.ReadyToStart{
		NextGameStart: t.Add(time.Duration(s.NextGameRemainingTicks/tickRate) * time.Second).Unix(),
	})
	log.Printf("buf:%v", buf)
	log.Printf("err:%v", err)
}
