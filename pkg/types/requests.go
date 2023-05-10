package types

import (
	"encoding/json"
	"fmt"
)

type InitResponse struct {
	Resp InitResponseFormat `json:"response"`
}

type InitResponseFormat struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	Ts     string `json:"ts"`
}

type InitRequest struct {
	Token   string
	GroupID string
	V       string
}

type WaitUpdatesRequest struct {
	Key string
	Ts  string
}

type WaitUpdatesResponse struct {
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
	Action Action `json:"action"`
}

type Action struct {
	Type string `json:"type"`
}

type ClientInfo struct {
	ButtonActions  []string `json:"button_actions"`
	Keyboard       bool     `json:"keyboard"`
	InlineKeyboard bool     `json:"inline_keyboard"`
	Carousel       bool     `json:"carousel"`
	LangID         int      `json:"lang_id"`
}

type SendMessageRequest struct {
	Token    string
	UserID   string
	Random   string
	Text     string
	Keyboard string
	V        string
}

type SendMessageResponse struct {
	Resp int
}

type ResponseError struct {
	err     error
	content string
}

type ErrorResponse struct {
	Error *VKError
}

type VKError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (e *VKError) Error() error {
	return fmt.Errorf(fmt.Sprint(e.ErrorCode) + " " + e.ErrorMsg)
}

type Button struct {
	Action struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
		Label   string `json:"label"`
	} `json:"action"`
	Color string `json:"color"`
}

type Keyboard struct {
	Inline  bool       `json:"inline"`
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
}

func (kbd Keyboard) Bytes() []byte {
	b, _ := json.Marshal(kbd)
	return b
}

type Reply struct {
	Msg      string
	Keyboard *Keyboard
}

func NewButton(label string, payload interface{}) Button {
	button := Button{}
	button.Action.Type = "text"
	button.Action.Label = label
	button.Action.Payload = "{}"
	if payload != nil {
		jPayoad, err := json.Marshal(payload)
		if err == nil {
			button.Action.Payload = string(jPayoad)
		}
	}
	button.Color = "default"
	return button
}

func NewKeyboard(buttons ...Button) *Keyboard {
	keyboard := Keyboard{Buttons: make([][]Button, 0)}
	row := make([]Button, 0)
	for _, b := range buttons {
		row = append(row, b)
	}
	keyboard.Buttons = append(keyboard.Buttons, row)
	return &keyboard
}

func (kbd *Keyboard) Append(b Button) {
	row := make([]Button, 0)
	row = append(row, b)
	kbd.Buttons = append(kbd.Buttons, row)
}
