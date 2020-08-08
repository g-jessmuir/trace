package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"time"

	"github.com/g-jessmuir/trace/trace"
)

func main() {
	fmt.Println("Tracing")
	start := time.Now()
	// get seed value somehow
	img := trace.GoTrace(99, 8)
	elapsed := time.Since(start)
	fmt.Println("time elapsed:", elapsed)

	f, err := os.Create("client/img.jpg")
	if err != nil {
		panic(err)
	}

	if err = jpeg.Encode(f, img, nil); err != nil {
		f.Close()
		panic(err)
	}
}
