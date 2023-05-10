package bot

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/qwertyqq2/vk-chat-testtask/configs"
	"github.com/qwertyqq2/vk-chat-testtask/pkg/types"
)

func loadConif(path string) (configs.Config, error) {
	conf, err := configs.LoadConfig(path)
	if err != nil {
		return configs.Config{}, err
	}
	return conf, nil
}

func procTesting(t *testing.T, name string, conf configs.Config, proc func(bot *Bot) error) {
	mybot := NewBot(
		GroupID(conf.GroupID),
		Token(conf.Token),
		Debug(true),
	)
	if err := mybot.Init(); err != nil {
		t.Fatal(err)
	}
	if err := proc(mybot); err != nil {
		t.Fatal(err)
	}
}

var conf configs.Config

func init() {
	path := "../../configs/envs"
	var err error
	conf, err = loadConif(path)
	if err != nil {
		log.Fatal(err)
	}

}

func TestCallbacks(t *testing.T) {
	procTesting(t, "build", conf, func(bot *Bot) error {
		name1, name2, name3 := "1", "2", "3"
		names := []string{name1, name2, name3}
		begin := NewInitCallback(name1)

		item1 := NewCallback(name2, "1")
		item2 := NewCallback(name3, "2")

		begin.AddNext(item1, item2)

		bot.Build(begin)

		getting := []string{}
		for k, _ := range bot.callbacks {
			getting = append(getting, k)
		}
		if len(getting) != len(names) {
			return fmt.Errorf("incorrect callbacks store")
		}
		return nil
	})

}

func TestSend(t *testing.T) {
	procTesting(t, "sendGroup", conf, func(bot *Bot) error {
		group, err := strconv.Atoi(conf.GroupID)
		if err != nil {
			return err
		}
		if err := bot.Send(group, types.DefualtResponse); err != nil {
			return err
		}
		return nil
	})

	procTesting(t, "sendMe", conf, func(bot *Bot) error {
		me, err := strconv.Atoi(conf.UserID)
		if err != nil {
			return err
		}
		if err := bot.Send(me, types.DefualtResponse); err != nil {
			return err
		}
		return nil
	})
}
