package trace

import (
	"image"
	"image/color"
	"math"
)

func getColor(r Ray, hl HitList) Vec {
	var rec HitRecord
	if hl.Hit(r, 0.0, math.MaxFloat32, &rec) {
		return Vec{rec.N.X + 1, rec.N.Y + 1, rec.N.Z + 1}.Mul(0.5)
	}
	u := r.Dir().Unit()
	t := 0.5 * (u.Y + 1.0)
	left := Vec{1.0, 1.0, 1.0}.Mul(1.0 - t)
	right := Vec{0.5, 0.7, 1.0}.Mul(t)
	return left.Add(right)
}

// Trace returns a pointer to the traced image
func Trace() *image.NRGBA {
	nx := 200
	ny := 100
	img := image.NewNRGBA(image.Rect(0, 0, nx, ny))
	lowerLeft := Vec{-2, -1, -1}
	horiz := Vec{4, 0, 0}
	vert := Vec{0, 2, 0}
	origin := Vec{0, 0, 0}
	hl := HitList{
		Sphere{Vec{0, 0, -1}, 0.5},
		Sphere{Vec{0, -100.5, -1}, 100},
	}
	fctob := func(f float32) byte { return byte(255.99 * f) }
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			u := float32(i) / float32(nx)
			v := float32(j) / float32(ny)
			r := Ray{origin, lowerLeft.Add(horiz.Mul(u)).Add(vert.Mul(v))}
			col := getColor(r, hl)
			ir := fctob(col.X)
			ig := fctob(col.Y)
			ib := fctob(col.Z)
			img.Set(i, ny-j-1, color.NRGBA{ir, ig, ib, 255})
		}
	}

	return img
}
