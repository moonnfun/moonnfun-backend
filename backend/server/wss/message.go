package wss

import (
	"fmt"
)

const C_Msg_Chat = "chat"

const C_FROM_SYSTEM = "system"

/*
init, send: {"type": "connection_init"}, response: {"type": "connection_ack"}
*/
type WMsg struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Wait    bool   `json:"-"`
	Payload any    `json:"payload"`
}

func WrapMsg(clientID, msgID, t string, payload any) *WMsg {
	if payload == nil {
		return &WMsg{Type: t} // fmt.Sprintf(`{"type": "%s"}`, t)
	}
	wmsg := &WMsg{ID: fmt.Sprintf("%s_%s", clientID, msgID), Type: t, Payload: payload}
	// wbuf, _ := json.Marshal(wmsg)
	// return string(wbuf)
	return wmsg
}
