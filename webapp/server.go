package webapp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func wsHandler(res http.ResponseWriter, req *http.Request) error {
	c, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		return fmt.Errorf("error upgrading request %v, error %v", req, err)
	}
	defer c.Close()

	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			return fmt.Errorf("error reading message, error %v", err)
		}

		log.Printf("msg type %d msg %q", msgType, msg)

		if err := c.WriteMessage(msgType, []byte("ok")); err != nil {
			return fmt.Errorf("error writing message, error %v", err)
		}
	}
}

// ListenAndServe starts serving the webapp
func ListenAndServe() {
	handleBowerDist := func(name string) {
		h := http.StripPrefix(
			"/"+name+"/",
			FileServerNoReaddir(http.Dir("webapp/bower_components/"+name+"/dist")),
		)

		http.Handle("/"+name+"/", h)
	}

	// bower components
	handleBowerDist("bootstrap")
	handleBowerDist("jquery")
	handleBowerDist("bacon")

	// local assets
	http.Handle("/js/", FileServerNoReaddir(http.Dir("webapp")))
	http.Handle("/", FileServerNoReaddir(http.Dir("webapp/html")))

	// endpoints
	http.Handle("/ws", HandlerIgnoreErr(DecorateHandler(HandlerErrFunc(wsHandler), logger)))

	http.ListenAndServe(":8080", nil)
}
