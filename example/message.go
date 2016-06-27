package main

import (
	"encoding/json"
)

//Message is the top envelop for message communication between nodes
type Message struct {
	ID int `json:"id"`
}

// CreateMessage creates a new message and returns a pointer
func CreateMessage() *Message {
	return &Message{}
}

// Serialize returns a byte array of the message
func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// Set function updates the content of message m awwording to input n
// And it fills the Msg's interface Contract
func (m *Message) Set(n []byte) error {
	type input struct {
		ID int `json:"int"`
	}
	var message input
	err := json.Unmarshal(n, &message)
	if err != nil {
		return err
	}
	m.ID = message.ID
	return nil
}
