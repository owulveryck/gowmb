package gowmb

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

// CreateHandler takes an "OutMessage" creator and returns a handler
func CreateHandler(creator OutMessage) func(http.ResponseWriter, *http.Request) {

	// ServeWs handles websocket requests from the peer.
	return func(w http.ResponseWriter, r *http.Request) {
		//Let's get the Tag
		vars := mux.Vars(r)
		Tag, err := stringToTag(vars["tag"])
		if err != nil {
			log.Warn("No Tag provided, bailing out")
			return
		}
		var contextLogger = log.WithFields(log.Fields{
			"Tag":  Tag,
			"From": r.RemoteAddr,
		})
		contextLogger.Info("New connection")
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		conn := &Conn{send: make(chan OutMessage, 256), ws: ws}
		reply := &reply{
			Message: creator,
			Tag:     Tag,
			Rep:     make(chan *hub),
		}
		defer close(reply.Rep)

		AllHubs.Request <- reply
		hub := <-reply.Rep
		hub.register <- conn
		go conn.writePump()
		conn.readPump(hub)
		contextLogger.Info("Connection ended")
	}

}
