// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wss

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
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

type AppendMsg struct {
	Msg    *WMsg
	Ktime  string
	SendCh chan struct{}
	WaitCh chan struct{}
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

	Active bool

	Topic    string
	WaitInit bool

	WaitKlines sync.Map
}

func NewClient(id string, conn *websocket.Conn) *Client {
	client := &Client{ID: id, conn: conn, TimeoutCh: time.After(pongWait)}

	// // wait read
	go client.waitRecv(context.Background())

	// conn.CloseRead(context.Background())

	return client
}

// A -> B, system -> B
func (c *Client) Push(msgID, msgType string, payload any, bInit bool) {
	if !bInit {
		if !c.Active && msgType != C_Msg_connection_init {
			return
		}
	}
	if len(strings.Split(c.Topic, "-")) == 2 {
		taddress := strings.Split(c.Topic, "-")[0]
		topicTime := strings.Split(c.Topic, "-")[1]
		if taddress != "" && strings.Contains(msgID, taddress) {
			if topicTime != "" && strings.HasSuffix(msgType, topicTime) {
				c.doPush(msgID, msgType, payload, bInit)
			}
		}
	}
}

func (c *Client) doPush(msgID, msgType string, payload any, bInit bool) {
	if msgID == "" {
		msgID = fmt.Sprintf("system_%v", time.Now().UnixNano())
	}
	msg := WrapMsg(c.ID, msgID, msgType, payload)

	// // 为了处理1分钟，5分钟，4小时等换挡临界态，存在数据乱序导致图表异常的问题，
	// // 采取延迟推送Append消息，直到第一次出现同等时间段的Update消息
	// strTimestamp := strings.Split(msgID, "_")[1]
	// if strings.HasPrefix(msgType, global.C_Msg_Price_Append) {
	// 	appendMsg := &AppendMsg{
	// 		Msg:    msg,
	// 		Ktime:  strTimestamp,
	// 		SendCh: make(chan struct{}),
	// 		WaitCh: make(chan struct{}),
	// 	}
	// 	c.WaitKlines.Store(strTimestamp, appendMsg)
	// 	go func(appendMsg *AppendMsg) {
	// 		for {
	// 			select {
	// 			case <-appendMsg.SendCh:
	// 				appendMsg.SendCh = nil
	// 				if err := c.sendMessage(context.Background(), appendMsg.Msg); err != nil {
	// 					slog.Error("send append message failed", "msg", msg, "error", err.Error())
	// 				}

	// 				c.WaitKlines.Delete(appendMsg.Ktime)
	// 				close(appendMsg.WaitCh)
	// 				return
	// 			case <-time.After(10 * time.Second):
	// 				global.Debug("send append message successed with timeout 10s", "msg", msg)
	// 				close(appendMsg.SendCh)
	// 				return
	// 			}
	// 		}
	// 	}(appendMsg)
	// 	return
	// } else if strings.HasPrefix(msgType, global.C_Msg_Price_Update) {
	// 	if appendMsg, ok := c.WaitKlines.Load(strTimestamp); ok {
	// 		if appendMsg.(*AppendMsg).SendCh != nil {
	// 			global.Debug("append kline", "timestamp", strTimestamp)
	// 			close(appendMsg.(*AppendMsg).SendCh)
	// 		}
	// 		<-appendMsg.(*AppendMsg).WaitCh
	// 	}
	// 	if strings.HasSuffix(msgType, "minute_1") {
	// 		kline := payload.([]model.TokenOHLCV)[0]
	// 		global.Debug("update kline", "timestamp", strTimestamp, "UnixMilli", time.UnixMilli(int64(kline.OT)).UnixMilli())
	// 	}
	// }

	// push msg to client
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
	// global.Debug("before send msg to client", "c.Topic", c.Topic, "message", wmsg)
	if c.conn != nil {
		if c.WaitInit && !strings.Contains(wmsg.Type, "init") {
			return nil
		}
		c.WaitInit = false
		if err := wsjson.Write(ctx, c.conn, wmsg); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid connection: %s", c.ID)
	}
	return nil
}
