package bot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qwertyqq2/vk-chat-testtask/pkg/types"
	"github.com/qwertyqq2/vk-chat-testtask/pkg/utils"
)

const (
	defaultUpdTimerInterval = 2 * time.Second

	initMessage = "Начать"
	backMessage = "Назад"
)

type Bot struct {
	Config

	ctx    context.Context
	cancel func()

	callbacks map[string]*Callback
	backUsers map[int]*backEvent
	hlk       sync.RWMutex

	upd          chan types.UpdMessage
	updTimer     time.Timer
	updInteraval time.Duration

	debug bool
}

type backEvent struct {
	prev, cur string
}

func NewBot(opts ...Option) *Bot {
	var conf Config
	for _, o := range opts {
		o(&conf)
	}
	return &Bot{
		Config:       conf,
		callbacks:    make(map[string]*Callback),
		backUsers:    make(map[int]*backEvent),
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
			if upd.Object.Message.Action.Type == "chat_invite_user" {
				bot.debugMessage("invite user")
				text = initMessage
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

		case "message_event":

		default:
			bot.debugMessage("undefined message type")
		}
	}
}

func (bot *Bot) sendCallback(sender int, event string, text string) {
	fmt.Println("event: ", sender, event, text)
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
				} else {
					bot.updateBackUser(u.Sender, u.Text)
				}
			}
		}
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

func (bot *Bot) execHandler(name string, userID int) bool {
	bot.hlk.RLock()
	if name == backMessage {
		e, ok := bot.backUsers[userID]
		if !ok {
			bot.updateBackUser(userID, initMessage)
		}
		prevHand, _ := bot.callbacks[e.prev]
		bot.hlk.RUnlock()

		if err := prevHand.handler(userID); err != nil {
			bot.debugMessage("err handler back " + err.Error())
			return false
		}
		bot.updateBackUser(userID, e.prev)
		return true
	} else {
		cb, ok := bot.callbacks[name]
		if !ok {
			bot.debugMessage("err handler callback ")
			return false
		}
		bot.hlk.RUnlock()

		if err := cb.handler(userID); err != nil {
			bot.debugMessage("err handler " + err.Error())
			return false
		}
		bot.updateBackUser(userID, name)

		return true
	}
}

func (bot *Bot) updateBackUser(userID int, name string) {
	bot.hlk.Lock()
	defer bot.hlk.Unlock()
	e, ok := bot.backUsers[userID]
	if !ok {
		e = &backEvent{prev: name, cur: name}
		bot.backUsers[userID] = e
	}
	if e.cur == name {
		return
	}
	call, ok := bot.callbacks[name]
	if !ok {
		return
	}
	e.cur = name
	if call.prev == nil {
		return
	}
	e.prev = call.prev.name
}

type Callback struct {
	handler func(int) error
	name    string
	prev    *Callback
	next    []*Callback
	message string
	opt     string
}

func NewCallback(name, text, opt string) *Callback {
	return &Callback{
		name:    name,
		next:    make([]*Callback, 0),
		message: text,
		opt:     opt,
	}
}

func NewInitCallback(text string) *Callback {
	return &Callback{
		name:    initMessage,
		next:    make([]*Callback, 0),
		message: text,
	}
}

func NewCallbackWithMessage(name, mes, opt string) *Callback {
	return &Callback{
		name:    name,
		message: mes,
		opt:     opt,
	}
}

func (c *Callback) AddNext(others ...*Callback) {
	for _, other := range others {
		if other.prev != nil {
			return
		}
		c.next = append(c.next, other)
		other.prev = c
	}
}

func (bot *Bot) Build(c *Callback) {
	buttons := make([]types.Button, 0, len(c.next))
	if c.next != nil {
		for _, call := range c.next {
			color := call.opt
			if color == "" {
				color = "primary"
			}
			b := types.NewButton(call.name, nil, false, color)
			buttons = append(buttons, b)
		}
	}
	if c.prev != nil {
		back := types.NewButton(backMessage, nil, false, "")
		buttons = append(buttons, back)
	}
	c.handler = func(userID int) error {
		kbd := types.NewKeyboard(buttons...).Bytes()
		var mes string
		mes = c.name
		if c.message != "" {
			mes = c.message
		}
		return bot.Send(userID, types.MessageSend{
			Text:     mes,
			Keyboard: string(kbd),
		})
	}
	bot.callbacks[c.name] = c
	for _, call := range c.next {
		bot.Build(call)
	}
}

func (bot *Bot) debugMessage(str string) {
	if bot.debug == true {
		fmt.Println(str)
	}
}
