package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

type H map[string]string

const (
	vkAPIURL        = "https://api.vk.com/method/"
	vkAPIVer        = "5.131"
	messagesCount   = 200
	requestInterval = 400 // 3 requests per second VK limit
	longPollVersion = 3

	availableCommands = "Available commands: /help, /me"
)

const (
	groupID = "220399914"
	token   = "vk1.a.VRqelsETm87AjI91mnRV7oZwuKOz3VL4EWSOu9Osi3iULVVhTy9bdtW89HvfHIF871oJyqPpm6t3GCavpcNrQk2b0GB2fo9yzDKGslpAHV0BhQifDbnUodfvkdCt7UZAzP-p8nsAI2r_2HTKSsjxb8HmAJdn1Fb9OjoHcn5kjMDm-Z2-BRvZz0i-u1mWfeBM4hozhgI4JVQS2Eo1g1AIiw"
	key     = "86944e7afac4f0ae0a8149371d2426a52cb13be2"
	server  = "https://lp.vk.com/wh220399914"
	ts      = "33"
)

type Config struct {
	groupID string
	token   string
	server  string
	key     string
	ts      string
}

var defaultConfig = Config{
	groupID: groupID,
	token:   token,
	server:  server,
	key:     key,
	ts:      ts,
}

type APIHandler struct {
	Config
}

func newDefaultAPIHandler() *APIHandler {
	return &APIHandler{
		Config: defaultConfig,
	}
}

func New() *APIHandler {
	return newDefaultAPIHandler()
}

func (api *APIHandler) Init() error {
	params := make(map[string]string)
	params["access_token"] = api.token
	params["group_id"] = api.groupID
	params["v"] = vkAPIVer
	parameters := url.Values{}
	for k, v := range params {
		parameters.Add(k, v)
	}

	resp, err := http.PostForm(vkAPIURL+"groups.getLongPollServer", parameters)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var res InitResponse
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return err
	}
	api.key = res.Resp.Key
	api.server = res.Resp.Server
	api.ts = res.Resp.Ts
	return nil
}

func (api *APIHandler) GetUpdate() error {
	params := make(map[string]string)
	params["act"] = "a_check"
	params["key"] = api.key
	params["ts"] = api.ts
	params["wait"] = "25"

	parameters := url.Values{}
	for k, v := range params {
		parameters.Add(k, v)
	}

	resp, err := http.PostForm(api.server, parameters)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, buf, "", "\t"); err != nil {
		return err
	}

	fmt.Println(string(prettyJSON.Bytes()))

	var res LongResponse
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

func (api *APIHandler) Send(text string) error {
	params := make(map[string]string)
	params["access_token"] = api.token
	params["user_id"] = "476713092"
	params["random_id"] = randID()
	params["message"] = text
	params["v"] = vkAPIVer
	parameters := url.Values{}
	for k, v := range params {
		parameters.Add(k, v)
	}

	resp, err := http.PostForm(vkAPIURL+"/messages.send", parameters)
	if err != nil {
		return fmt.Errorf("post err " + err.Error())
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := ErrorResponse{}
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return fmt.Errorf("vkapi: vk response is not json: " + string(buf))
	}

	if r.Error != nil {
		return r.Error.strVkError()
	}

	var res interface{}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return fmt.Errorf("unmarshal " + err.Error())
	}
	fmt.Println(res)
	return nil
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
	//	RequestParams
}

func (e *VKError) strVkError() error {
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

// Keyboard to send for user
type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
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

func setKeyboard() Keyboard {
	keyboard := Keyboard{Buttons: make([][]Button, 0)}
	button1 := NewButton("/me", nil)
	button2 := NewButton("/create", nil)
	row := make([]Button, 0)
	row = append(row, button1)
	row = append(row, button2)
	keyboard.Buttons = append(keyboard.Buttons, row)
	return keyboard
}

func (api *APIHandler) SendKeyboard() error {
	keyb := setKeyboard()
	keybJson, err := json.Marshal(keyb)
	if err != nil {
		return err
	}
	fmt.Println(string(keybJson))
	params := make(map[string]string)
	params["access_token"] = api.token
	params["user_id"] = "476713092"
	params["random_id"] = randID()
	params["message"] = "some"
	params["keyboard"] = string(keybJson)
	params["v"] = vkAPIVer
	parameters := url.Values{}
	for k, v := range params {
		parameters.Add(k, v)
	}

	resp, err := http.PostForm(vkAPIURL+"/messages.send", parameters)
	if err != nil {
		return fmt.Errorf("post err " + err.Error())
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := ErrorResponse{}
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return fmt.Errorf("vkapi: vk response is not json: " + string(buf))
	}

	if r.Error != nil {
		return r.Error.strVkError()
	}

	var res interface{}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return fmt.Errorf("unmarshal " + err.Error())
	}
	fmt.Println(res)
	return nil
}

// https://lp.vk.com/wh220399914?act=a_check&key=184ccb9c705fec1fbcf5c1f4fc57e486f3847f14&ts=33&wait=25

func randID() string {
	return strconv.FormatUint(uint64(rand.Uint32()), 10)
}
