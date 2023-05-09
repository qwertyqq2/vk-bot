package api

type InitResponseFormat struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	Ts     string `json:"ts"`
}

type InitResponse struct {
	Resp InitResponseFormat `json:"response"`
}

type LongResponse struct {
	Ts      string   `json:"ts"`
	Updates []Update `json:"updates"`
}

type Update struct {
	Type    string `json:"type"`
	EventID string `json:"event_id"`
	Object  Object `json:"object"`
	GroupID int    `json:"group_id"`
	V       string `json:"v"`
}

type Object struct {
	Message    Message    `json:"message"`
	ClientInfo ClientInfo `json:"client_info"`
}

type Message struct {
	ID     int    `json:"id"`
	Date   int    `json:"date"`
	PeerID int    `json:"peer_id"`
	FromID int    `json:"from_id"`
	Text   string `json:"text"`
}

type ClientInfo struct {
	ButtonActions  []string `json:"button_actions"`
	Keyboard       bool     `json:"keyboard"`
	InlineKeyboard bool     `json:"inline_keyboard"`
	Carousel       bool     `json:"carousel"`
	LangID         int      `json:"lang_id"`
}

type ButtonAction []string
