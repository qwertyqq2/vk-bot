package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/qwertyqq2/vk-chat-testtask/pkg/types"
)

const (
	reqGetLong      = "groups.getLongPollServer"
	baseUrl         = "https://api.vk.com/method/"
	version         = "5.131"
	initAddr        = "groups.getLongPollServer"
	messageSendAddr = "/messages.send"
)

type RequestFormat map[string]string

func InitRequest(req types.InitRequest) (types.InitResponse, error) {
	var resp types.InitResponse
	err := createRequest(RequestFormat{
		"access_token": req.Token,
		"group_id":     req.GroupID,
		"v":            req.V,
	}, &resp, baseUrl+initAddr)
	if err != nil {
		return types.InitResponse{}, err
	}
	return resp, nil
}

func WaitUpdatesRequest(req types.WaitUpdatesRequest, serverAddr string) (types.WaitUpdatesResponse, error) {
	var resp types.WaitUpdatesResponse
	err := createRequest(RequestFormat{
		"act":  "a_check",
		"key":  req.Key,
		"ts":   req.Ts,
		"wait": "25",
	}, &resp, serverAddr)
	if err != nil {
		return types.WaitUpdatesResponse{}, err
	}
	return resp, nil
}

func SendMessageRequest(req types.SendMessageRequest) (types.SendMessageResponse, error) {
	var resp types.SendMessageResponse
	err := createRequest(RequestFormat{
		"access_token": req.Token,
		"user_id":      req.UserID,
		"random_id":    req.Random,
		"message":      req.Text,
		"keyboard":     req.Keyboard,
		"v":            req.V,
	}, &resp, baseUrl+messageSendAddr)
	if err != nil {
		return types.SendMessageResponse{}, err
	}
	return resp, nil
}

func createRequest(req RequestFormat, val interface{}, addr string) error {
	parameters := url.Values{}
	for k, v := range req {
		parameters.Add(k, v)
	}
	resp, err := http.PostForm(addr, parameters)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := types.ErrorResponse{}
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return fmt.Errorf("vkapi: vk response is not json: " + string(buf))
	}
	if r.Error != nil {
		return r.Error.Error()
	}

	err = json.Unmarshal(buf, val)
	if err != nil {
		return err
	}
	return nil
}
