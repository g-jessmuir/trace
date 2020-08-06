package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

type msg struct {
	Num int
}

func echo(conn *websocket.Conn) {
	for {
		m := msg{}

		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
		}

		fmt.Printf("Got message: %#v\n", m)
		if err = conn.WriteJSON("blablabla"); err != nil {
			fmt.Println(err)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go echo(conn)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	a := r.URL.Path
	fmt.Println(a)
	if a == "/" {
		a = "/index.html"
	}
	content, err := ioutil.ReadFile("client/" + a)
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}

func main() {
	// fmt.Println("Tracing")
	// start := time.Now()
	// img := trace.GoTrace(16, nil)
	// elapsed := time.Since(start)
	// fmt.Println("time elapsed:", elapsed)

	// f, err := os.Create("client/img.png")
	// if err != nil {
	// 	panic(err)
	// }

	// if err = png.Encode(f, img); err != nil {
	// 	f.Close()
	// 	panic(err)
	// }

	fmt.Println("serving!")
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", rootHandler)

	panic(http.ListenAndServe(":9090", nil))
}
