package main

import (
	"log"

	"github.com/qwertyqq2/vk-chat-testtask/configs"
	"github.com/qwertyqq2/vk-chat-testtask/data"
	"github.com/qwertyqq2/vk-chat-testtask/pkg/bot"
)

var pages *data.Data

func main() {
	if err := loadData(); err != nil {
		log.Fatal(err)
	}

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

	begin := bot.NewInitCallback("Здравствуте, на связи магазин компьютерных программ)")

	item1 := bot.NewCallback("Программы", "Выберите тип программы")
	item2 := bot.NewCallback("Игры", "Выберете тип игры")
	item3 := bot.NewCallback("Подписки", "Выберите тип подписки?")
	item4 := bot.NewCallback("Контакты", "Как вы предпочитаете связаться с нами?")

	begin.AddNext(item1, item2, item3, item4)

	item11 := bot.NewCallbackWithMessage("Платные", pages.Content("pay_progs"))
	item12 := bot.NewCallbackWithMessage("Бесплатные", pages.Content("unpay_progs"))

	item21 := bot.NewCallbackWithMessage("Одиночные", pages.Content("only_games"))
	item22 := bot.NewCallbackWithMessage("Сетевые", pages.Content("multy_games"))

	item31 := bot.NewCallbackWithMessage("Фильмы", pages.Content("subs_movie"))
	item32 := bot.NewCallbackWithMessage("Сериалы", pages.Content("subs_series"))

	item41 := bot.NewCallbackWithMessage("Связь", pages.Content("conn_msg"))
	item42 := bot.NewCallbackWithMessage("Оставить отзыв", pages.Content("conn_comm"))

	item1.AddNext(item11, item12)
	item2.AddNext(item21, item22)
	item3.AddNext(item31, item32)
	item4.AddNext(item41, item42)

	mybot.Build(begin)

	select {}
}

func loadData() error {
	d, err := data.NewData("conn_comm", "conn_msg", "multy_games",
		"only_games", "pay_progs", "subs_movie", "subs_series", "unpay_progs")
	if err != nil {
		return err
	}
	pages = d
	return nil
}
