package types

type UpdMessage struct {
	Sender   int
	Message  string
	Text     string
	Date     int
	Keyboard bool
}

type MessageSend struct {
	Text     string
	Keyboard string
}

var (
	DefaultUndefinedNameMes = MessageSend{Text: "undefined name of object"}
	DefualtResponse         = MessageSend{Text: "ok"}
)
