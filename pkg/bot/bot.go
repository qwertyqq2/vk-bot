package bot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qwertyqq2/vk-chat-testtask/pkg/session"
	"github.com/qwertyqq2/vk-chat-testtask/pkg/types"
	"github.com/qwertyqq2/vk-chat-testtask/pkg/utils"
)

const (
	defaultUpdTimerInterval = 2 * time.Second
)

type callback struct {
	handler func(int) error
	name    string
	prev    string
}

func newcallback(name, prev string, handler func(int) error) callback {
	return callback{handler: handler, name: name, prev: prev}
}

type Bot struct {
	Config

	ctx    context.Context
	cancel func()

	sessions map[string]*session.Session
	slk      sync.RWMutex

	//layers       map[string]*Layer
	//handleLayers map[string]func(int) error
	callbacks map[string]callback
	hlk       sync.RWMutex

	upd          chan types.UpdMessage
	updTimer     time.Timer
	updInteraval time.Duration

	debug bool
}

func NewBot(opts ...Option) *Bot {
	var conf Config
	for _, o := range opts {
		o(&conf)
	}
	return &Bot{
		Config:       conf,
		sessions:     make(map[string]*session.Session),
		callbacks:    make(map[string]callback),
		upd:          make(chan types.UpdMessage, 100),
		updTimer:     *time.NewTimer(defaultUpdTimerInterval),
		updInteraval: defaultUpdTimerInterval,
		debug:        conf.debug,
	}
}

func (bot *Bot) Init() error {
	if bot.token == "" {
		return fmt.Errorf("undefined token")
	}
	if bot.groupID == "" {
		return fmt.Errorf("undefined groupID")
	}
	if bot.v == "" {
		bot.debugMessage("verison dont set, setting default")
		bot.v = version
	}
	resp, err := InitRequest(types.InitRequest{
		Token:   bot.token,
		GroupID: bot.groupID,
		V:       bot.v,
	})
	if err != nil {
		return err
	}

	if resp.Resp.Key == "" || resp.Resp.Server == "" || resp.Resp.Ts == "" {
		return fmt.Errorf("incorrect api init response")
	}

	ctx, cancel := context.WithCancel(context.Background())

	bot.key = resp.Resp.Key
	bot.server = resp.Resp.Server
	bot.ts = resp.Resp.Ts
	bot.ctx = ctx
	bot.cancel = cancel

	go bot.checkUpdates()
	go bot.run()
	return nil
}

func (bot *Bot) checkUpdates() {
	tiker := time.NewTicker(defaultUpdTimerInterval)

	for {
		select {
		case <-bot.ctx.Done():
			return
		case <-tiker.C:
			upds, err := WaitUpdatesRequest(types.WaitUpdatesRequest{
				Key: bot.key,
				Ts:  bot.ts,
			}, bot.server)
			if err != nil {
				bot.debugMessage("err wait response")
				continue
			}
			bot.ts = upds.Ts
			go bot.handleUpds(upds)
		}
	}
}

func (bot *Bot) handleUpds(upds types.WaitUpdatesResponse) {
	for _, upd := range upds.Updates {
		switch upd.Type {
		case "message_new":
			sender := upd.Object.Message.FromID
			text := upd.Object.Message.Text
			time := upd.Object.Message.Date
			if sender == 0 || time == 0 || text == "" {
				bot.debugMessage(fmt.Sprintf("incorrect message"))
				continue
			}
			select {
			case <-bot.ctx.Done():
				return
			default:
			}

			go func(keyboard bool) {
				select {
				case <-bot.ctx.Done():
					return
				case bot.upd <- types.UpdMessage{
					Sender:   sender,
					Text:     text,
					Date:     time,
					Keyboard: keyboard,
				}:
				}
			}(upd.Object.ClientInfo.Keyboard)

		default:
			bot.debugMessage("undefined message type")
		}
	}
}

func (bot *Bot) run() {
	for {
		select {
		case <-bot.ctx.Done():
			return
		case u := <-bot.upd:
			bot.debugMessage(fmt.Sprintf("receive mes %d, %s, %d, %t", u.Sender, u.Text, u.Date, u.Keyboard))
			if u.Keyboard {
				if !bot.execHandler(u.Text, u.Sender) {
					err := bot.Send(u.Sender, types.MessageSend{
						Text: "undefined obj name",
					})
					if err != nil {
						bot.debugMessage("cant send message")
					}
				}
			}
		}
	}
}

func (bot *Bot) ShowCallbacks() {
	for k, _ := range bot.callbacks {
		fmt.Println(k)
	}
}

func (bot *Bot) Send(userID int, mes types.MessageSend) error {
	_, err := SendMessageRequest(types.SendMessageRequest{
		Token:    bot.token,
		UserID:   fmt.Sprint(userID),
		Random:   utils.RandID(),
		Text:     mes.Text,
		Keyboard: mes.Keyboard,
		V:        bot.v,
	})
	return err
}

func (bot *Bot) Callback(name, prev string, handler func(int) error) types.Button {
	button := types.NewButton(name, nil)

	bot.hlk.Lock()
	defer bot.hlk.Unlock()

	cb := newcallback(name, prev, handler)
	bot.callbacks[cb.name] = cb
	return button
}

func (bot *Bot) NewKeyboard(prev string, buttons ...types.Button) *types.Keyboard {
	kbd := types.NewKeyboard(buttons...)
	if prev != "" {
		back := bot.Callback("back", prev, func(userID int) error {
			if prev == "" {
				return nil
			}
			cb, ok := bot.callbacks[prev]
			if !ok {
				return fmt.Errorf("undefined prev callback")
			}
			if cb.prev == "" {
				return nil
			}
			next, ok := bot.callbacks[cb.prev]
			if !ok {
				return fmt.Errorf("undefined parent callback")
			}
			if err := next.handler(userID); err != nil {
				return err
			}
			return nil
		})
		kbd.Append(back)
	}
	return kbd
}

func (bot *Bot) execHandler(name string, userID int) bool {
	bot.hlk.RLock()
	cb, ok := bot.callbacks[name]
	if !ok {
		return false
	}
	bot.hlk.RUnlock()

	if err := cb.handler(userID); err != nil {
		bot.debugMessage("err handler " + err.Error())
	}
	return true
}

func (bot *Bot) debugMessage(str string) {
	if bot.debug == true {
		fmt.Println(str)
	}
}
