// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wss

import (
	"context"
	"fmt"
	"log/slog"
	"meme3/global"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	C_Msg_connection_init = "init"
	C_Msg_connection_ack  = "ack"
	C_Msg_topic_subscribe = "subscribe"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 15 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

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

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID string

	PushCh chan struct{}

	TimeoutCh <-chan time.Time

	// The websocket connection.
	conn *websocket.Conn

	// // Buffered channel of outbound messages.
	// send chan []byte

	Active bool

	Topics []string
}

func NewClient(id string, conn *websocket.Conn) *Client {
	client := &Client{ID: id, conn: conn, TimeoutCh: time.After(pongWait)}
	client.Topics = make([]string, 0)

	// // wait read
	go client.waitRecv(context.Background())

	// conn.CloseRead(context.Background())

	return client
}

// A -> B, system -> B
func (c *Client) Push(msgID, msgType string, payload any, bInit bool) {
	global.Debug("client push", slog.Any("client", c))
	if !bInit {
		if !c.Active && msgType != C_Msg_connection_init {
			return
		}
	}

	if msgID == "" {
		msgID = fmt.Sprintf("system_%v", time.Now().UnixNano())
	}
	msg := WrapMsg(c.ID, msgID, msgType, payload)

	// wbuf, _ := json.Marshal(msg)
	// global.Debug("before send msg to client", "msgID", msgID, "msg", string(wbuf))
	// c.send <- wbuf

	global.Debug("before send msg to client", "msg", msg)
	if err := c.sendMessage(context.Background(), msg); err != nil {
		slog.Error("send message failed", "msg", msg, "error", err.Error())
	}
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close(websocket.StatusInternalError, "internal error")
		c.conn = nil
	}
}

func (c *Client) waitRecv(ctx context.Context) {
	defer c.Close()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if c.conn == nil {
				ticker.Stop()
				return
			}
			_, msgData, err := c.conn.Read(ctx)
			if err != nil {
				// slog.Error("read client message failed with UnexpectedClosedError", "error", err.Error())
				// return
				continue
			}
			if string(msgData) == "o" {
				c.TimeoutCh = time.After(pongWait)
				go c.SendPing(ctx)
				continue
			}
			go MsgHandle(c, msgData)
		case <-c.TimeoutCh:
			slog.Error("read client message failed with UnexpectedClosedError", "error", "timeout")
			ticker.Stop()
			return
		}
	}
}

func (c *Client) SendPing(ctx context.Context) {
	<-time.After(time.Duration(pingPeriod))
	if c.conn != nil {
		c.conn.Write(ctx, websocket.MessageText, []byte("p"))
		// if err := c.conn.Ping(ctx); err != nil {
		// 	slog.Error("send ping failed", "error", err.Error())
		// }
	}
}

func (c *Client) sendMessage(ctx context.Context, wmsg *WMsg) error {
	global.DebugForce("before send msg to client", "c.Topics", c.Topics, "message", wmsg)
	if c.conn != nil {
		if len(c.Topics) > 0 && strings.Contains(wmsg.Type, c.Topics[0]) {
			if err := wsjson.Write(ctx, c.conn, wmsg); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("invalid connection: %s", c.ID)
	}
	return nil
}
