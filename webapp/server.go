package webapp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/s-urbaniak/rec/msg"
	"github.com/s-urbaniak/rec/recorder"
)

var upgrader = websocket.Upgrader{}

func wsHandler(res http.ResponseWriter, req *http.Request) error {
	c, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		return fmt.Errorf("error upgrading request %v, error %v", req, err)
	}
	defer c.Close()

	var rec recorder.Recorder
	var msgChan chan msg.Msg
	for {
		msgType, m, err := c.ReadMessage()
		if err != nil {
			return fmt.Errorf("error reading message, error %v", err)
		}

		log.Printf("msgType %d msg %q", msgType, m)
		cmd := string(m)

		switch {
		case cmd == "3" && rec == nil:
			log.Println("recording start")
			rec = recorder.NewRecorder()
			rec.Start()
			msgChan = make(chan msg.Msg)
			rec.MsgChan(msgChan)

			go func() {
				println("monitoring")
				for v := range msgChan {
					switch m := v.(type) {
					case msg.MsgLevel:
						t := (m.Timestamp - (m.Timestamp % 1000000)) / 1000000
						millis := t % 1000
						t = t / 1000 // sec
						sec := t % 60
						t = t / 60 // min
						ts := fmt.Sprintf("%02d:%02d:%03d", t, sec, millis)

						if err := c.WriteMessage(msgType, []byte(ts)); err != nil {
							panic(fmt.Sprintf("error writing message, error %v", err))
						}
					}
				}
				println("done monitoring")
			}()

		case cmd == "3" && rec != nil:
			log.Println("recording stop")
			rec.Stop()
			close(msgChan)
			rec = nil
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
