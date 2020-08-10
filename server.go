package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/g-jessmuir/trace/trace"
	"github.com/gorilla/websocket"
)

type message struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func traceHandler(conn *websocket.Conn) {
	for {
		var gottenMsg trace.Settings
		err := conn.ReadJSON(&gottenMsg)
		if err != nil {
			fmt.Println("Error reading json:", err)
			return
		}

		fmt.Println("Got seed:", gottenMsg.Seed)
		fmt.Println("Tracing")
		start := time.Now()
		updater := make(chan string, 3) // add some buffer for latency
		endSignal := make(chan int)
		// get seed value somehow
		var imgString string
		go trace.GoTrace(gottenMsg, updater, endSignal)
		go func() {
			for {
				select {
				case imgString = <-updater:
					err := conn.WriteJSON(message{"working", imgString})
					if err != nil {
						fmt.Println("error writing:", err)
					}
				case <-endSignal:
					elapsed := time.Since(start)
					fmt.Println("time elapsed:", elapsed)
					err := conn.WriteJSON(message{"done", imgString})
					if err != nil {
						fmt.Println("error writing:", err)
					}
					return
				}
			}
		}()
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	go traceHandler(conn)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	path := "client/" + r.URL.Path[1:]
	if path == "client/" {
		path = "client/index.html"
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not open file:", err)
	}
	fmt.Fprintf(w, "%s", content)
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", rootHandler)
	fmt.Println("serving!")
	panic(http.ListenAndServe(":8080", nil))
}
