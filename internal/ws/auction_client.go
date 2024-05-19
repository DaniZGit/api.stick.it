package ws

import (
	"fmt"
	"log"
	"time"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/gorilla/websocket"
)

const (
	AuctionEventTypeBid = "auction_event_bid"
	AuctionEventTypeCompleted = "auction_event_completed"
	AuctionEventTypeCreated = "auction_event_created"
)

type AuctionEvent struct {
	Type string `json:"type"`
	Payload any `json:"payload"`
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readAuctionPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		fmt.Println("New incoming message", string(message))
		// c.hub.Broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writeAuctionPump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	
	for {
		select {
			case message, ok := <-c.send:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// The hub closed the channel.
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(message)

				// Add queued chat messages to the current websocket message.
				n := len(c.send)
				for i := 0; i < n; i++ {
					w.Write(<-c.send)
				}

				if err := w.Close(); err != nil {
					return
				}
			case <-ticker.C:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeAuctionWs(hub *Hub, ctx *app.ApiContext) {
	// create ws connection
	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	
	// register client to the hub
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeAuctionPump()
	go client.readAuctionPump()
}