package message

type Message struct {
	ID          int
	Date        int
	Out         int
	UserID      int   `json:"user_id"`
	ChatID      int   `json:"chat_id"`
	PeerID      int64 `json:"peer_id"`
	ReadState   int   `json:"read_state"`
	Title       string
	Body        string
	Action      string
	ActionMID   int `json:"action_mid"`
	Flags       int
	Timestamp   int64
	Payload     string
	FwdMessages []Message `json:"fwd_messages"`
}

// Messages - VK Messages
type Messages struct {
	Count int
	Items []*Message
}

type MessagesResponse struct {
	Response Messages
	Error    *VKError
}

type VKError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	//	RequestParams
}
