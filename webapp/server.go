package webapp

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func wsHandler(res http.ResponseWriter, req *http.Request) {
	c, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Printf("error upgrading request %v, error %v", req, err)
		return
	}
	defer c.Close()

	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("error reading message, error %v", err)
			break
		}
		log.Printf("msg type %d msg %q", msgType, msg)

		if err := c.WriteMessage(msgType, []byte("ok")); err != nil {
			log.Printf("error writing message, error %v", err)
			break
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
	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe(":8080", nil)
}
