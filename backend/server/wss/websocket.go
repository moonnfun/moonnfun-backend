package wss

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"meme3/global"
	"strings"
	"sync"

	// // "luckyton/server/
	// "meme3/service/logic/core"
	// "meme3/service/logic/funds"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients sync.Map

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	global.WebsocketHandler = HttpToWebsocket
}

func HttpToWebsocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	if c, ok := clients.Load(conn.RemoteAddr().String()); ok && c != nil {
		slog.Info("Before reset websocket client", "websocketClient", c)
		c.(*Client).Close()
		// return fmt.Errorf("connect websocket again and again, address: %s", conn.RemoteAddr().String())
	}

	// add client
	client := NewClient(conn.RemoteAddr().String(), conn)
	clients.Store(conn.RemoteAddr().String(), client)
	slog.Info("Add websocket client", "address", conn.RemoteAddr().String())

	// // player, err := logic.GetUser(api.WebParams(r).Get("address"), true)
	// // if err != nil {
	// // 	slog.Error("WebsocketHandler.PlayerVerify", "error", err.Error())
	// // 	conn.Close()
	// // 	// api.WebResponseJson(w, r, api.ApiError(err.Error()), http.StatusInternalServerError)
	// // 	return err
	// // }

	// // force reset client
	// if player.WebsocketClient != nil {
	// 	slog.Info("Before reset websocket client", "player", player.Address, "websocketClient", player.WebsocketClient)
	// 	player.WebsocketClient = nil
	// }

	// // add client
	// client := NewClient(player.Address, conn)
	// slog.Info("Add websocket client", "player", player.Address, "address", r.RemoteAddr)
	// // whub.AddClient(client)
	// player.WebsocketClient = client

	// // balance update push
	// listHotFunds := funds.GetListFunds(player.TonFunds.Hot)
	// global.Debug("before push funds update", "player", player.ID, "listHotFunds", listHotFunds)
	// go client.Push("", global.C_Msg_funds_update, listHotFunds, true)

	// fundsnotify := funds.GetFundsNotify("Channel_CryptoBot", "Ton", "1", decimal.NewFromInt(1), false)
	// slog.Info("push deposited notify successed", "player", player.ID, "token", "Ton", "deposited", "1")
	// go client.Push("", global.C_Msg_notify_push, fundsnotify, true)
	return nil
}

func MsgHandle(client *Client, msg string) (bHandle bool, retErr error) {
	wmsg := &WMsg{}
	if err := json.Unmarshal([]byte(msg), wmsg); err != nil {
		return false, err
	}
	slog.Info("receive websocket msg successed", "msg", wmsg)

	if wmsg.Type == C_Msg_connection_init {
		client.Active = true
		client.Push("", C_Msg_connection_ack, nil, false)
	} else if wmsg.Type == C_Msg_topic_subscribe {
		// format: address-second_1
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
		client.PushCh = global.WebsocketSubscribe(client.conn.RemoteAddr().String(), address, topic)
	}
	return nil
}

func WebsocketSend(clientID, tokenAddress, msgType string, payload any) error {
	client, ok := clients.Load(clientID)
	if ok && client != nil {
		global.Debug("before send msg to client", "clientID", clientID, "tokenAddress", tokenAddress, "msgType", msgType, "payload", payload)
		client.(*Client).Push(tokenAddress, msgType, payload, false)
		return nil
	} else {
		clients.Range(func(key, value any) bool {
			if c, ok := value.(*Client); ok && c != nil && len(c.Topics) > 0 && strings.HasSuffix(msgType, c.Topics[0]) {
				global.Debug("before broadcast msg to client", "toID", key, "msgType", msgType, "payload", payload)
				value.(*Client).Push(tokenAddress, msgType, payload, false)
			}
			return true
		})
	}
	return nil

	// if toID != "" {
	// 	if player, err := logic.GetUser(toID, true); player != nil {
	// 		player.WebsocketClient.(*Client).Push(fromID, msgType, payload, false)
	// 		return nil
	// 	} else {
	// 		return fmt.Errorf("websocket subscribe failed with invalid id: %s, from: %s, err: %s", toID, fromID, err.Error())
	// 	}
	// } else {
	// 	global.Debug("before braodcast msg", "from", fromID, "msgType", msgType)
	// 	return errors.New("unsupport broadcast")
	// }
}
