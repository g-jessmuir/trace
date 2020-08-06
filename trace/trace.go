package trace

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

func randomScene() HitList {
	n := 500
	list := make(HitList, n+1)
	list[0] = Sphere{Vec{0, -1000, 0}, 1000, Lambertian{Vec{0.5, 0.5, 0.5}}}
	index := 1
	for i := -11; i < 11; i++ {
		for j := -11; j < 11; j++ {
			chooseMat := rand.Float32()
			center := Vec{float32(i) + 0.9*rand.Float32(), 0.2, float32(j) + 0.9*rand.Float32()}
			if center.Sub(Vec{4, 0.2, 0}).Len() > 0.9 {
				if chooseMat < 0.8 {
					// diffuse
					col := Vec{
						rand.Float32() * rand.Float32(),
						rand.Float32() * rand.Float32(),
						rand.Float32() * rand.Float32(),
					}
					list[index] = Sphere{center, 0.2, Lambertian{col}}
					index++
				} else if chooseMat < 0.95 {
					// metal
					col := Vec{
						0.5 * (1 + rand.Float32()),
						0.5 * (1 + rand.Float32()),
						0.5 * (1 + rand.Float32()),
					}
					list[index] = Sphere{center, 0.2, Metal{col, 0.3 * rand.Float32()}}
					index++
				} else {
					// glass
					list[index] = Sphere{center, 0.2, Dielec{1.5}}
					index++
				}
			}
		}
	}
	list[index] = Sphere{Vec{0, 1, 0}, 1, Dielec{1.5}}
	index++
	list[index] = Sphere{Vec{-4, 1, 0}, 1, Lambertian{Vec{0.2, 0.5, 0.3}}}
	index++
	list[index] = Sphere{Vec{4, 1, 0}, 1, Metal{Vec{0.7, 0.6, 0.5}, 0.0}}
	index++
	return list[:index]
}

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
	lookFrom := Vec{12, 2, 5}
	lookAt := Vec{0, 0.5, 0}
	distToFocus := lookFrom.Sub(lookAt).Len()
	cam := CreateCam(20, float32(nx)/float32(ny), lookFrom, lookAt, Vec{0, 1, 0}, 0.2, distToFocus)
	// hl := HitList{
	// 	Sphere{Vec{0, 0, -1}, 0.5, Lambertian{Vec{0.1, 0.2, 0.5}}},
	// 	Sphere{Vec{0, -100.5, -1}, 100, Lambertian{Vec{0.8, 0.8, 0.0}}},
	// 	Sphere{Vec{1, 0, -1}, 0.5, Metal{Vec{0.8, 0.6, 0.2}, 0.3}},
	// 	Sphere{Vec{-1, 0, -1}, 0.5, Dielec{1.5}},
	// 	Sphere{Vec{-1, 0, -1}, -0.45, Dielec{1.5}},
	// }

	hl := randomScene()

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
