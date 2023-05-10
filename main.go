package main

import (
	"log"

	"github.com/qwertyqq2/vk-chat-testtask/pkg/bot"
	"github.com/qwertyqq2/vk-chat-testtask/pkg/types"
)

func main() {
	mybot := bot.NewBot(
		bot.GroupID("220399914"),
		bot.Token("vk1.a.VRqelsETm87AjI91mnRV7oZwuKOz3VL4EWSOu9Osi3iULVVhTy9bdtW89HvfHIF871oJyqPpm6t3GCavpcNrQk2b0GB2fo9yzDKGslpAHV0BhQifDbnUodfvkdCt7UZAzP-p8nsAI2r_2HTKSsjxb8HmAJdn1Fb9OjoHcn5kjMDm-Z2-BRvZz0i-u1mWfeBM4hozhgI4JVQS2Eo1g1AIiw"),
		bot.Debug(true),
	)

	if err := mybot.Init(); err != nil {
		log.Fatal(err)
	}

	begin := bot.NewCallback("begin")

	item1 := bot.NewCallback("item1")
	item2 := bot.NewCallback("item2")

	begin.AddNext(item1)
	begin.AddNext(item2)

	item11 := bot.NewCallbackWithHander("item11", func(userID int) error {
		return mybot.Send(userID, types.DefualtResponse)
	})
	item21 := bot.NewCallbackWithHander("item21", func(userID int) error {
		return mybot.Send(userID, types.DefualtResponse)
	})
	item1.AddNext(item11)
	item2.AddNext(item21)

	mybot.Build(begin)

	select {}
}
