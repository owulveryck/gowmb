package gowmb_test

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/owulveryck/gowmb"
	"time"

	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	testServer *httptest.Server
	reader     io.Reader //Ignore this for now
	baseWsURL  string
)

func init() {
	router := mux.NewRouter().StrictSlash(true)

	handler := gowmb.CreateHandler(createMessage(), newTag(), "tag")
	router.
		Methods("GET").
		Path("/serveWs/{tag}").
		Name("WebSocket").
		HandlerFunc(handler)

	testServer = httptest.NewServer(router) //Creating new server with the user handlers

	baseWsURL = fmt.Sprintf("%s/serveWs/", testServer.URL) //Grab the address for the endpoint

}

func TestPingPong(t *testing.T) {
	tsURL, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}
	wsURL := url.URL{Scheme: "ws", Host: tsURL.Host, Path: "/serveWs/1234"}
	c, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Errorf("Cannot connect to the websocket %v", err)

	}
	defer c.Close()
	if err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(2*time.Second)); err != nil {
		t.Errorf("write close: %v", err)
	}
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, ""))
	if err != nil {
		t.Errorf("write close: %v", err)
	}

}
func TestServeWs(t *testing.T) {
	tsURL, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}
	httpURL := url.URL{Scheme: tsURL.Scheme, Host: tsURL.Host, Path: "/serveWs/"}
	// Try to connect to a socket without an ID
	request, err := http.NewRequest("GET", httpURL.String(), nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// We don't serve the baseurl, a tag is mandatory
	if res.StatusCode != 404 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

	//Try with a valid tag
	httpURL.Path = "/serveWs/1234'"
	request, err = http.NewRequest("GET", httpURL.String(), nil)

	res, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// We shall get a bad request as we are expected a websocket
	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	// Now test the websocket
	wsURL := url.URL{Scheme: "ws", Host: tsURL.Host, Path: "/serveWs/1234"}
	c, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Errorf("Cannot connect to the websocket %v", err)

	}
	defer c.Close()

	done := make(chan bool)

	go func() {
		defer close(done)
		_, message, err := c.ReadMessage()
		if err != nil {
			t.Errorf("read: %v", err)
		}
		t.Logf("recv: %s", message)
		done <- true
	}()
	type input struct {
		ID int `json:"id"`
	}

	message := &input{ID: 0}
	b, err := json.Marshal(message)
	if err != nil {
		t.Error(err)
	}
	//reader = strings.NewReader(string(b))
	t.Logf("Writing...")
	err = c.WriteMessage(websocket.TextMessage, b)
	t.Logf("Done...")
	if err != nil {
		t.Errorf("write: %v", err)
	}
	t.Logf("Waiting...")
	<-done
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, ""))
	if err != nil {
		t.Errorf("write close: %v", err)
	}
	t.Logf("End...")
	//<-time.After(5 * time.Second)
}
