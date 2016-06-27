package gowmb_test

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/owulveryck/gowmb"
	"log"
	"net/http"
	"strconv"
)

//message is the top envelop for message communication between nodes
type message struct {
	ID int `json:"id"`
}

// Createmessage creates a new message and returns a pointer
func createMessage() *message {
	return &message{}
}

// Serialize returns a byte array of the message
func (m *message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// Set function updates the content of message m awwording to input n
// And it fills the Msg's interface Contract
func (m *message) Set(n []byte) error {
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

type tag int

func (t *tag) Parse(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*t = tag(v)
	return nil
}

func newTag() *tag {
	return new(tag)
}

func Example() {
	router := mux.NewRouter().StrictSlash(true)

	handler := gowmb.CreateHandler(createMessage(), newTag(), "tag")
	router.
		Methods("GET").
		Path("/serveWs/{tag}").
		Name("WebSocket").
		HandlerFunc(handler)
	log.Fatal(http.ListenAndServe(":8080", router))

}
