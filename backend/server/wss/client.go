// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wss

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	C_Msg_connection_init = "init"
	C_Msg_connection_ack  = "ack"
	C_Msg_topic_subscribe = "subscribe"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID string

	PushCh chan struct{}

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	Active bool

	Topics []string
}

func NewClient(id string, conn *websocket.Conn) *Client {
	client := &Client{ID: id, conn: conn, send: make(chan []byte, 256)}
	client.Topics = make([]string, 0)
	go client.readPump()
	go client.writePump()
	return client
}

// A -> B, system -> B
func (c *Client) Push(fromID, msgType string, payload any, bInit bool) {
	slog.Debug("client push", slog.Any("client", c))
	if !bInit {
		if !c.Active && msgType != C_Msg_connection_init {
			return
		}
	}

	if fromID == "" {
		fromID = "system"
	}
	msg := WrapMsg(fromID, c.ID, msgType, payload)
	wbuf, _ := json.Marshal(msg)
	slog.Debug("before send msg to client", "toID", c.ID, "msg", string(wbuf))
	c.send <- wbuf
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		// c.hub.unregister <- c
		c.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("read client message failed with UnexpectedClosedError", "error", err.Error())
			} else {
				slog.Error("read client message failed with ClosedError", "error", err.Error())
			}
			break
		}
		slog.Debug("Client.readPump receive msg successed", "msg", message)
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		go MsgHandle(c, string(message))
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		// c.conn.Close()
		c.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				slog.Warn("read message from client's send channel failed", "id", c.ID)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.Active = false
				return
			}

			slog.Debug("Client.writePump receive msg successed", "msg", message)
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Error("generate client's writer failed", "msg", message)
				return
			}
			w.Write(message)
			slog.Debug("system send msg to websocket client successed", "msg", message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("send msg to client with Ping failed", "error", err.Error())
				return
			}
		}
	}
}
