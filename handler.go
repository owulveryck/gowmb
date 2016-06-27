package gowmb

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/owulveryck/topology-presentation/message"
	"net/http"
)

// ServeWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
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
		Message: message.CreateMessage(),
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
