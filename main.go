package main

import (
	"log"

	"github.com/qwertyqq2/vk-chat-testtask/pkg/bot"
	"github.com/qwertyqq2/vk-chat-testtask/pkg/types"
)

func main() {
	bot := bot.NewBot(
		bot.GroupID("220399914"),
		bot.Token("vk1.a.VRqelsETm87AjI91mnRV7oZwuKOz3VL4EWSOu9Osi3iULVVhTy9bdtW89HvfHIF871oJyqPpm6t3GCavpcNrQk2b0GB2fo9yzDKGslpAHV0BhQifDbnUodfvkdCt7UZAzP-p8nsAI2r_2HTKSsjxb8HmAJdn1Fb9OjoHcn5kjMDm-Z2-BRvZz0i-u1mWfeBM4hozhgI4JVQS2Eo1g1AIiw"),
		bot.Debug(true),
	)

	if err := bot.Init(); err != nil {
		log.Fatal(err)
	}

	bot.Callback("/begin", "", func(userID int) error {
		item1 := bot.Callback("/item1", "/begin", func(userID int) error {
			item1 := bot.Callback("item11", "/item1", func(userID int) error {
				return bot.Send(userID, types.MessageSend{
					Text: "item11",
				})
			})
			item2 := bot.Callback("item12", "/item1", func(userID int) error {
				return bot.Send(userID, types.MessageSend{
					Text: "item12",
				})
			})
			kbd := bot.NewKeyboard("/item1", item1, item2).Bytes()
			return bot.Send(userID, types.MessageSend{
				Text:     "second lawer",
				Keyboard: string(kbd),
			})
		})

		item2 := bot.Callback("/item2", "/begin", func(userID int) error {
			item1 := bot.Callback("item11", "/item2", func(userID int) error {
				return bot.Send(userID, types.MessageSend{
					Text: "item21",
				})
			})
			item2 := bot.Callback("item12", "/item2", func(userID int) error {
				return bot.Send(userID, types.MessageSend{
					Text: "item22",
				})
			})
			kbd := bot.NewKeyboard("/item2", item1, item2).Bytes()
			return bot.Send(userID, types.MessageSend{
				Text:     "second lawer",
				Keyboard: string(kbd),
			})
		})

		item3 := bot.Callback("/item3", "/begin", func(userID int) error {
			item1 := bot.Callback("item11", "/item3", func(userID int) error {
				return bot.Send(userID, types.MessageSend{
					Text: "item31",
				})
			})
			item2 := bot.Callback("item12", "/item3", func(userID int) error {
				return bot.Send(userID, types.MessageSend{
					Text: "item32",
				})
			})
			kbd := bot.NewKeyboard("/item3", item1, item2).Bytes()
			return bot.Send(userID, types.MessageSend{
				Text:     "second lawer",
				Keyboard: string(kbd),
			})
		})

		item4 := bot.Callback("/item4", "/begin", func(userID int) error {
			item1 := bot.Callback("item41", "/item4", func(userID int) error {
				return bot.Send(userID, types.MessageSend{
					Text: "item41",
				})
			})
			item2 := bot.Callback("item42", "/item4", func(userID int) error {
				return bot.Send(userID, types.MessageSend{
					Text: "item42",
				})
			})
			kbd := bot.NewKeyboard("/item4", item1, item2).Bytes()
			return bot.Send(userID, types.MessageSend{
				Text:     "second lawer",
				Keyboard: string(kbd),
			})
		})
		kbd := bot.NewKeyboard("", item1, item2, item3, item4).Bytes()
		return bot.Send(userID, types.MessageSend{
			Text:     "begin",
			Keyboard: string(kbd),
		})
	})

	bot.ShowCallbacks()

	select {}
}
