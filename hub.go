// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gowmb

import (
	log "github.com/Sirupsen/logrus"
)

func init() {
	go allHubs.Run()
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Tag of the hubs
	Tag Tag

	// Registered connections.
	connections map[*Conn]bool

	// The last message broadcasted
	message *Messager

	// Inbound messages from the connections for processing purpose.
	process chan []byte

	// Inbound messages from the connections.
	broadcast chan Messager

	// Register requests from the connections.
	register chan *Conn

	// Unregister requests from connections.
	unregister chan *Conn
}

// hubs maintains the set of active hubs
type hubs struct {
	// Unregister
	unregister chan Tag

	// Registered hubs.
	hubs map[Tag]*hub

	// Register requests from the connections.
	Request chan *reply
}

type reply struct {
	Tag     Tag
	Message Messager
	Rep     chan *hub
}

// allHubs is the actual registry of all hubs
var allHubs = hubs{
	unregister: make(chan Tag),
	Request:    make(chan *reply),
	hubs:       make(map[Tag]*hub),
}

// Run is main routine for registering the hubs
func (h *hubs) Run() {
	for {
		select {
		case r := <-h.Request:
			if _, ok := h.hubs[r.Tag]; !ok {
				//TODO create a new hub
				var hub = &hub{
					Tag:         r.Tag,
					message:     &r.Message,
					process:     make(chan []byte),
					broadcast:   make(chan Messager),
					register:    make(chan *Conn),
					unregister:  make(chan *Conn),
					connections: make(map[*Conn]bool),
				}
				var contextLogger = log.WithFields(log.Fields{
					"Tag": r.Tag,
					"hub": &hub,
				})
				contextLogger.Debug("New HUB")
				go hub.run()
				h.hubs[r.Tag] = hub
				// By the end reply to the sender
			}
			r.Rep <- h.hubs[r.Tag]
		case hub := <-h.unregister:
			log.Debug("In the hubs' unregister")
			if _, ok := h.hubs[hub]; ok {
				var contextLogger = log.WithFields(log.Fields{
					"Tag": hub,
					"hub": h.hubs[hub],
				})
				contextLogger.Debug("Unregistering HUB")
				delete(h.hubs, hub)
			}
		}
	}
}
func (h *hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn] = true
			log.WithFields(log.Fields{
				"Connections": len(h.connections),
				"Connection":  conn,
				"hub":         &h,
			}).Debug("Registerng connection")
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				log.WithFields(log.Fields{
					"Connections": len(h.connections),
					"Connection":  conn,
					"hub":         &h,
				}).Debug("Unregisterng connection")
				delete(h.connections, conn)
				close(conn.send)
			}
			// If the last element has been removed exit)
			if len(h.connections) == 0 {
				allHubs.unregister <- h.Tag
				return
			}
		case inMessage := <-h.process:
			log.WithFields(log.Fields{
				"hub": &h,
			}).Debug("process")
			err := (*h.message).Set(inMessage)
			if err != nil {
				log.WithFields(log.Fields{
					"hub": &h,
				}).Debug("Cannot process message, diescarding it", err)
			} else {
				go func() {
					h.broadcast <- *h.message
				}()
			}
		case outMessage := <-h.broadcast:
			log.WithFields(log.Fields{
				"hub": &h,
			}).Debug("Broadcast")
			for conn := range h.connections {
				log.WithFields(log.Fields{
					"Connection": conn,
					"hub":        &h,
				}).Debug("Sending...")
				select {
				case conn.send <- outMessage:
				default:
					close(conn.send)
					//delete(hub.connections, conn)
					delete(h.connections, conn)
				}
			}
		}
	}
}
