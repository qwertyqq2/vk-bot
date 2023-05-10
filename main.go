package main

import (
	"log"

	"github.com/qwertyqq2/vk-chat-testtask/configs"
	"github.com/qwertyqq2/vk-chat-testtask/pkg/bot"
)

func main() {
	conf, err := configs.LoadConfig("./configs/envs")
	if err != nil {
		log.Fatal(err)
	}
	mybot := bot.NewBot(
		bot.GroupID(conf.GroupID),
		bot.Token(conf.Token),
		bot.Debug(true),
	)

	if err := mybot.Init(); err != nil {
		log.Fatal(err)
	}

	begin := bot.NewCallback("begin")

	item1 := bot.NewCallback("Программы")
	item2 := bot.NewCallback("Игры")
	item3 := bot.NewCallback("Подписки")
	item4 := bot.NewCallback("Контакты")

	begin.AddNext(item1, item2, item3, item4)

	item11 := bot.NewCallbackWithMessage("Платные", "here1")
	item12 := bot.NewCallbackWithMessage("Бесплатные", "here1")

	item21 := bot.NewCallbackWithMessage("Одиночные", "here2")
	item22 := bot.NewCallbackWithMessage("Сетевые", "here2")

	item31 := bot.NewCallbackWithMessage("Фильмы", "here3")
	item32 := bot.NewCallbackWithMessage("Сериалы", "here3")

	item41 := bot.NewCallbackWithMessage("Связь", "here4")
	item42 := bot.NewCallbackWithMessage("Оставить отзыв", "here4")

	item1.AddNext(item11, item12)
	item2.AddNext(item21, item22)
	item3.AddNext(item31, item32)
	item4.AddNext(item41, item42)

	mybot.Build(begin)

	select {}
}
