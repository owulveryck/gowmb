package gowmb

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

// CreateHandler returns a http.handler.
// It takes as input a Messager interface
func CreateHandler(creator Messager, tag Tag, tagName string) func(http.ResponseWriter, *http.Request) {

	// ServeWs handles websocket requests from the peer.
	return func(w http.ResponseWriter, r *http.Request) {
		//Let's get the Tag
		vars := mux.Vars(r)
		if _, ok := vars[tagName]; !ok {
			log.Errorf("Expected tag %v not found", tagName)
			return
		}
		err := tag.Parse(vars[tagName])
		if err != nil {
			log.Warn("Cannot parse tag %v: %v", vars[tagName], err)
			return
		}
		var contextLogger = log.WithFields(log.Fields{
			"Tag":  tag,
			"From": r.RemoteAddr,
		})
		contextLogger.Info("New connection")
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		conn := &Conn{send: make(chan Messager, 256), ws: ws}
		reply := &reply{
			Message: creator,
			Tag:     tag,
			Rep:     make(chan *hub),
		}
		defer close(reply.Rep)

		allHubs.Request <- reply
		hub := <-reply.Rep
		hub.register <- conn
		go conn.writePump()
		conn.readPump(hub)
		contextLogger.Info("Connection ended")
	}

}
