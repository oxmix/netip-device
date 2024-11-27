package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"netip-device/collector"
	"os"
	"os/signal"
	"time"
)

type connectResponse struct {
	responseBase
}

type connectPayload struct {
	payloadBase
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	destroy := make(chan bool, 1)

	connect(&connectPayload{
		payloadBase: payloadBase{
			Service: "device",
		},
	})

	col := collector.New()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// reader from nodes-handler
	go func() {
		for {
			if wsConnect == nil {
				time.Sleep(time.Second)
				continue
			}

			var res struct {
				Command string `json:"command"`
			}
			// ReadMessage needs for ping-pong
			err := wsConnect.ReadJSON(&res)
			if err != nil {
				if !websocket.IsCloseError(err, 1006) {
					log.Println("ws read err:", err)
				}
				time.Sleep(time.Second)
				continue
			}

			// debug
			// log.Println("received command:", res.Command)

			switch res.Command {
			case "services-destroy":
				destroy <- true
				return
			}
		}
	}()

	log.Println("ready to work")

	// writer to nodes-handler
	for {
		select {
		case <-destroy:
			log.Println("service destroyed")
			log.Println("------")
			log.Println("below remains to execution manually:")
			log.Println("# docker rm -f netip.device")
			log.Println("------")
			_ = wsConnect.Close()
			<-destroy

		// ticker-sender stats network
		case <-ticker.C:
			if wsConnect == nil {
				continue
			}
			j := struct {
				Event         string                  `json:"event"`
				CollectDevice collector.CollectDevice `json:"collectDevice"`
			}{
				Event:         "collect-device",
				CollectDevice: col.Get(),
			}
			err := wsConnect.WriteJSON(j)
			if err != nil {
				// ! lock this select before established connection
				connectDegrade(errors.New("ticker write err: " + err.Error()))
				continue
			}

		// chan-sender stats disks info
		case cdi, ok := <-col.ChanDisksInfo:
			if wsConnect == nil || !ok {
				continue
			}
			j := &struct {
				Event     string               `json:"event"`
				DisksInfo *collector.DisksInfo `json:"disksInfo"`
			}{
				Event:     "disks-info",
				DisksInfo: cdi,
			}
			err := wsConnect.WriteJSON(j)
			if err != nil {
				log.Println("disks info write err:", err)
				continue
			}

		// handler terminate
		case <-interrupt:
			log.Println("interrupt")
			if wsConnect == nil {
				return
			}

			err := wsConnect.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close err:", err)
				return
			}
			err = wsConnect.Close()
			if err != nil {
				log.Println("close conn err:", err)
			}
			return
		}
	}
}
