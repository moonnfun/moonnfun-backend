package wss

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/server/api"
	"meme3/service/logic"

	// // "luckyton/server/
	// "meme3/service/logic/core"
	// "meme3/service/logic/funds"
	"net/http"

	"github.com/gorilla/websocket"
)

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

	player, err := logic.GetUser(api.WebParams(r).Get("address"))
	if err != nil {
		slog.Error("WebsocketHandler.PlayerVerify", "error", err.Error())
		conn.Close()
		// api.WebResponseJson(w, r, api.ApiError(err.Error()), http.StatusInternalServerError)
		return err
	}

	// force reset client
	if player.WebsocketClient != nil {
		slog.Info("Before reset websocket client", "player", player.Address, "websocketClient", player.WebsocketClient)
		player.WebsocketClient = nil
	}

	// add client
	client := NewClient(player.Address, conn)
	slog.Info("Add websocket client", "player", player.Address, "address", r.RemoteAddr)
	// whub.AddClient(client)
	player.WebsocketClient = client

	// // balance update push
	// listHotFunds := funds.GetListFunds(player.TonFunds.Hot)
	// slog.Debug("before push funds update", "player", player.ID, "listHotFunds", listHotFunds)
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

	switch wmsg.Type {
	case C_Msg_connection_init:
		client.Active = true
		client.Push("", C_Msg_connection_ack, nil, false)
	case C_Msg_topic_subscribe:
		WebsocketSubscribe(client, wmsg.ID, fmt.Sprintf("%v", wmsg.Payload))
	// case C_Msg_game:
	// 	handleGameMsg(client, msg)
	// case C_Msg_notify:
	// 	handleGameMsg(client, msg)
	default:
		return false, nil
	}
	return true, nil
}

func WebsocketSubscribe(client *Client, id, topic string) error {
	if id == "" || topic == "" {
		return errors.New("topic subscribe failed with invalid params")
	}

	if client == nil {
		if player, err := logic.GetUser(id); player != nil {
			client = player.WebsocketClient.(*Client)
		} else {
			return fmt.Errorf("websocket subscribe failed with invalid id: %s, topic: %s, err: %s", id, topic, err.Error())
		}
	}
	slog.Info("before subscribe", "id", id, "topic", topic)
	client.Topics = append(client.Topics, topic)
	return nil
}

func WebsocketSend(fromID, toID, msgType string, payload any) error {
	if toID != "" {
		if player, err := logic.GetUser(toID); player != nil {
			player.WebsocketClient.(*Client).Push(fromID, msgType, payload, false)
			return nil
		} else {
			return fmt.Errorf("websocket subscribe failed with invalid id: %s, from: %s, err: %s", toID, fromID, err.Error())
		}
	} else {
		slog.Debug("before braodcast msg", "from", fromID, "msgType", msgType)
		return errors.New("unsupport broadcast")
	}
}
