package trace

import (
	"image"
	"image/color"
)

// Trace returns a pointer to the traced image
func Trace() *image.NRGBA {
	nx := 200
	ny := 100
	img := image.NewNRGBA(image.Rect(0, 0, nx, ny))
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			r := float32(i) / float32(nx)
			g := float32(j) / float32(ny)
			b := float32(0.2)
			fctob := func(f float32) byte { return byte(255.99 * f) }
			ir := fctob(r)
			ig := fctob(g)
			ib := fctob(b)
			img.Set(i, j, color.NRGBA{ir, ig, ib, 255})
		}
	}

	return img
}
