package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/g-jessmuir/trace/trace"
)

func main() {
	fmt.Println("Tracing")
	img := trace.Trace()

	f, err := os.Create("client/img.png")
	if err != nil {
		panic(err)
	}

	if err = png.Encode(f, img); err != nil {
		f.Close()
		panic(err)
	}

	fs := http.FileServer(http.Dir("client"))
	go log.Fatal(http.ListenAndServe(":9090", fs))
}
