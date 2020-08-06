package trace

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
)

// Args contains the settings for rendering
type Args struct {
	World  HitList
	Nx     int
	Ny     int
	Ns     int
	Camera Cam
}

type Pixel struct {
	r, g, b byte
	x, y    int
}

// randomScene creates a scene that mimics the
func randomScene() HitList {
	n := 500
	list := make(HitList, n)
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
func Trace(args Args, i int, j int, output chan Pixel) {
	fctob := func(f float32) byte { return byte(255.99 * f) }
	var col Vec
	for s := 0; s < args.Ns; s++ {
		u := (float32(i) + rand.Float32()) / float32(args.Nx)
		v := (float32(j) + rand.Float32()) / float32(args.Ny)
		r := args.Camera.getRay(u, v)
		col = col.Add(getColor(r, args.World, 0))
	}
	col = col.Div(float32(args.Ns))
	col = Vec{Sqrt32(col.X), Sqrt32(col.Y), Sqrt32(col.Z)}
	ir := fctob(col.X)
	ig := fctob(col.Y)
	ib := fctob(col.Z)
	output <- Pixel{r: ir, g: ig, b: ib, x: i, y: j}
}

func worker(args Args, input chan Pixel, output chan Pixel, id int) {
	for {
		select {
		case work := <-input:
			// fmt.Println("worker", id, " working on", work.x, work.y)
			Trace(args, work.x, work.y, output)
		default:
			return
		}
	}
}

// GoTrace accumulates multiple traces
func GoTrace(threads int, traceArgs *Args) *image.NRGBA {
	var args Args
	if traceArgs == nil {
		nx := 400
		ny := 200
		ns := 1
		lookFrom := Vec{12, 2, 5}
		lookAt := Vec{0, 0.5, 0}
		distToFocus := lookFrom.Sub(lookAt).Len()
		cam := CreateCam(15, float32(nx)/float32(ny), lookFrom, lookAt, Vec{0, 1, 0}, 0.2, distToFocus)
		args = Args{
			World:  randomScene(),
			Nx:     nx,
			Ny:     ny,
			Ns:     ns,
			Camera: cam,
		}
	} else {
		args = *traceArgs
	}
	work := make(chan Pixel, args.Nx*args.Ny)
	acc := make(chan Pixel, args.Nx*args.Ny)
	img := image.NewNRGBA(image.Rect(0, 0, args.Nx, args.Ny))
	for j := 0; j < args.Ny; j++ {
		for i := 0; i < args.Nx; i++ {
			work <- Pixel{x: i, y: j}
		}
	}
	for i := 0; i < threads; i++ {
		go worker(args, work, acc, i)
	}
	for c := 0; c < args.Nx*args.Ny; c++ {
		p := <-acc
		img.Set(p.x, args.Ny-p.y-1, color.NRGBA{p.r, p.g, p.b, 255})
		if c%(args.Nx*args.Ny/5) == 0 {
			fmt.Println("done%:", float32(c)/float32(args.Nx*args.Ny))
			// f, err := os.Create("client/img.png")
			// if err != nil {
			// 	panic(err)
			// }

			// if err = png.Encode(f, img); err != nil {
			// 	f.Close()
			// 	panic(err)
			// }
		}
	}
	return img
}
