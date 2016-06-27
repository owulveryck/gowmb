/*
Package gowmb (go Websocket Message Broker) manages message processing and broadcast
between various clients, based on a segregation "tag".

Basics

The CreateHandler should be used to create a http.Handler type that can be used with
the net/http (or any good http's muxer package).

Assuming you've defined a Messenger'c compatible message type, every client connected to
the web-path server by the defined handler will receive a Serialize version of the Message
each time the server's websocket detects an event.

*/
package gowmb
