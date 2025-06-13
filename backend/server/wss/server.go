package wss

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"meme3/global"
	"net/http"
	"strings"
	"sync"

	"github.com/coder/websocket"
)

var clients sync.Map

func Init(subscribeFunc func(id, address, topic string) chan struct{}) {
	global.WebsocketSend = WebsocketSend
	global.WebsocketHandler = HttpToWebsocket
	global.WebsocketHandler = HttpToWebsocket
	global.WebsocketSubscribe = subscribeFunc
}

func HttpToWebsocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}

	if c, ok := clients.Load(r.RemoteAddr); ok && c != nil {
		slog.Info("Before reset websocket client", "websocketClient", c)
		c.(*Client).Close()
		// return fmt.Errorf("connect websocket again and again, address: %s", conn.RemoteAddr().String())
	}

	// add client
	client := NewClient(r.RemoteAddr, conn)
	clients.Store(r.RemoteAddr, client)
	slog.Info("Add websocket client", "address", r.RemoteAddr)
	return nil
}

func MsgHandle(client *Client, msg []byte) (bHandle bool, retErr error) {
	wmsg := &WMsg{}
	if err := json.Unmarshal(msg, wmsg); err != nil {
		return false, err
	}

	if wmsg.Type == C_Msg_connection_init {
		client.Active = true
		client.Push("", C_Msg_connection_ack, nil, false)
		client.SendPing(context.Background())
		slog.Info("receive websocket msg successed", "msg", wmsg)
	} else if wmsg.Type == C_Msg_topic_subscribe {
		// format: address-second_1
		slog.Info("receive websocket msg successed", "msg", wmsg)
		subMsgs := strings.Split(fmt.Sprintf("%v", wmsg.Payload), "-")
		WebsocketSubscribe(client, wmsg.ID, subMsgs[0], subMsgs[1])
	} else {
		return false, nil
	}
	return true, nil
}

func WebsocketSubscribe(client *Client, id, address, topic string) error {
	if topic == "" {
		return errors.New("topic subscribe failed with invalid params")
	}

	if client == nil {
		return fmt.Errorf("websocket subscribe failed with invalid id: %s, topic: %s", id, topic)
	}
	slog.Info("before subscribe", "id", id, "topic", topic, "address", address)
	if len(client.Topics) == 1 {
		if client.PushCh != nil {
			close(client.PushCh)
			client.PushCh = nil
		}
		client.Topics[0] = topic
	} else {
		client.Topics = append(client.Topics, topic)
	}

	if global.WebsocketSubscribe != nil {
		client.PushCh = global.WebsocketSubscribe(client.ID, address, topic)
	}
	return nil
}

func WebsocketSend(clientID, tokenAddress, msgType string, payload any) error {
	client, ok := clients.Load(clientID)
	if ok && client != nil {
		global.Debug("before send msg to client", "clientID", clientID, "tokenAddress", tokenAddress, "msgType", msgType, "payload", payload)
		go client.(*Client).Push(tokenAddress, msgType, payload, false)
		return nil
	} else {
		clients.Range(func(key, value any) bool {
			if c, ok := value.(*Client); ok && c != nil && len(c.Topics) > 0 && strings.HasSuffix(msgType, c.Topics[0]) {
				global.Debug("before broadcast msg to client", "toID", key, "msgType", msgType, "payload", payload)
				go value.(*Client).Push(tokenAddress, msgType, payload, false)
			}
			return true
		})
	}
	return nil
}
