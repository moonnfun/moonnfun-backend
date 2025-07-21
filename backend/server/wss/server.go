package wss

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/service/store"
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

	// if c, ok := clients.Load(clientID); ok && c != nil {
	// 	slog.Info("Before reset websocket client", "websocketClient", c)
	// 	c.(*Client).Close()
	// 	// return fmt.Errorf("connect websocket again and again, address: %s", conn.RemoteAddr().String())
	// }
	clientId := store.NewId()
	clients.Store(clientId, NewClient(clientId, conn))
	slog.Info("Reveive websocket client", "clientId", clientId)

	// // add client
	// if clientID == "" {
	// 	clientID = store.NewId()
	// }
	// client := NewClient(clientID, conn)
	// clients.Store(clientID, client)
	// slog.Info("Add websocket client", "address", clientID)
	return nil
}

func MsgHandle(client *Client, msg []byte) (bHandle bool, retErr error) {
	wmsg := &WMsg{}
	if err := json.Unmarshal(msg, wmsg); err != nil {
		return false, err
	}

	if wmsg.Type == C_Msg_connection_init {
		client.Active = true
		oldID := fmt.Sprintf("%v", wmsg.Payload)
		go client.SendPing(context.Background())
		client.doPush("", C_Msg_connection_ack, client.ID, true)
		if oldID != client.ID && oldID != "" {
			if c, ok := clients.Load(oldID); ok && c != nil {
				slog.Info("Before reset websocket client", "websocketClient", c)
				WebsocketRemoveClient(oldID)
			}
		}
		slog.Info("receive websocket init successed", "msg", wmsg)
	} else if wmsg.Type == C_Msg_topic_subscribe {
		// format: address-second_1
		slog.Info("receive websocket subscribe successed", "msg", wmsg)
		subMsgs := strings.Split(fmt.Sprintf("%v", wmsg.Payload), "-")
		WebsocketSubscribe(client, client.ID, subMsgs[0], subMsgs[1])
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
		if c, _ := clients.Load(id); c == nil {
			return fmt.Errorf("websocket subscribe failed with invalid id: %s, topic: %s", id, topic)
		} else {
			client = c.(*Client)
			client.Active = true
		}
	}
	slog.Info("before subscribe", "client", client, "topic", topic, "address", address)
	client.Topic = fmt.Sprintf("%s-%s", address, topic)
	client.WaitInit = true

	if global.WebsocketSubscribe != nil {
		client.PushCh = global.WebsocketSubscribe(client.ID, address, topic)
	}
	return nil
}

func WebsocketSend(clientID, tokenAddress, msgType string, payload any) error {
	client, ok := clients.Load(clientID)
	if ok && client != nil {
		// global.Debug("before send msg to client", "clientID", clientID, "tokenAddress", tokenAddress, "msgType", msgType, "payload", payload)
		go client.(*Client).Push(tokenAddress, msgType, payload, false)
		return nil
	} else {
		clients.Range(func(key, value any) bool {
			if c, ok := value.(*Client); ok && c != nil && c.Topic != "" {
				// global.Debug("before broadcast msg to client", "toID", key, "msgType", msgType, "payload", payload)
				go value.(*Client).Push(tokenAddress, msgType, payload, false)
			}
			return true
		})
	}
	return nil
}

func WebsocketRemoveClient(clientID string) {
	if clientID != "" {
		clients.Delete(clientID)
	}
}

func WebsocketStatus(clientID string) bool {
	if c, ok := clients.Load(clientID); ok && c != nil && c.(*Client).GetConn() != nil {
		return true
	} else {
		return false
	}
}
