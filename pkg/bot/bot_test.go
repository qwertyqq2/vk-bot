package bot

import (
	"strconv"
	"testing"

	"github.com/qwertyqq2/vk-chat-testtask/configs"
	"github.com/qwertyqq2/vk-chat-testtask/pkg/types"
)

func loadConif(t *testing.T, path string) (configs.Config, error) {
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

func TestCallbacks(t *testing.T) {
	path := ""
	conf, err := loadConif(t, path)
	if err != nil {
		t.Fatal(err)
	}
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

	procTesting(t, "regCallback", conf, func(bot *Bot) error {
		// begin := NewInitCallback("")

		// item1 := bot.NewCallback("", "")
		// item2 := bot.NewCallback("Игры", "Выберете тип игры", "positive")
		return nil
	})

}
