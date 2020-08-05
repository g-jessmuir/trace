package trace

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

func getColor(r Ray, hl HitList, depth int) Vec {
	var rec HitRecord
	if hl.Hit(r, 0.001, math.MaxFloat32, &rec) {
		var scattered Ray
		var attenuation Vec
		if depth < 50 && rec.M.Scatter(r, &rec, &attenuation, &scattered) {
			return attenuation.VMul(getColor(scattered, hl, depth+1))
		}
		return Vec{0, 0, 0}
	}
	// Sky color
	u := r.Dir.Unit()
	t := 0.5 * (u.Y + 1.0)
	left := Vec{1.0, 1.0, 1.0}.Mul(1.0 - t)
	right := Vec{0.5, 0.7, 1.0}.Mul(t)
	return left.Add(right)
}

// Trace returns a pointer to the traced image
func Trace() *image.NRGBA {
	nx := 200
	ny := 100
	ns := 100
	img := image.NewNRGBA(image.Rect(0, 0, nx, ny))
	cam := Cam{
		Origin:     Vec{0, 0, 0},
		LowerLeft:  Vec{-2, -1, -1},
		Horizontal: Vec{4, 0, 0},
		Vertical:   Vec{0, 2, 0},
	}
	hl := HitList{
		Sphere{Vec{0, 0, -1}, 0.5, Lambertian{Vec{0.1, 0.2, 0.5}}},
		Sphere{Vec{0, -100.5, -1}, 100, Lambertian{Vec{0.8, 0.8, 0.0}}},
		Sphere{Vec{1, 0, -1}, 0.5, Metal{Vec{0.8, 0.6, 0.2}, 0.3}},
		Sphere{Vec{-1, 0, -1}, 0.5, Dielec{1.5}},
	}
	fctob := func(f float32) byte { return byte(255.99 * f) }
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			var col Vec
			for s := 0; s < ns; s++ {
				u := (float32(i) + rand.Float32()) / float32(nx)
				v := (float32(j) + rand.Float32()) / float32(ny)
				r := cam.GetRay(u, v)
				col = col.Add(getColor(r, hl, 0))
			}
			col = col.Div(float32(ns))
			col = Vec{Sqrt32(col.X), Sqrt32(col.Y), Sqrt32(col.Z)}
			ir := fctob(col.X)
			ig := fctob(col.Y)
			ib := fctob(col.Z)
			img.Set(i, ny-j-1, color.NRGBA{ir, ig, ib, 255})
		}
	}

	return img
}
