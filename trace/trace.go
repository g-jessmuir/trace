package trace

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"math/rand"
	"os"
	"time"
)

// Args contains the settings for rendering
type Args struct {
	World  HitList
	Nx     int
	Ny     int
	Ns     int
	Camera Cam
}

type pixel struct {
	r, g, b byte
	x, y    int
}

// randomScene creates a scene that mimics the scene in the book
// func randomScene() HitList {
// 	n := 500
// 	list := make(HitList, n)
// 	list[0] = Sphere{Vec{0, -500, 0}, 500, Lambertian{Vec{0.8, 0.8, 0.8}}}
// 	index := 1
// 	for i := -11; i < 11; i++ {
// 		for j := -11; j < 11; j++ {
// 			chooseMat := rand.Float32()
// 			center := Vec{float32(i) + 0.9*rand.Float32(), 0.2, float32(j) + 0.9*rand.Float32()}
// 			if center.Sub(Vec{4, 0.2, 0}).Len() > 0.9 {
// 				if chooseMat < 0.8 {
// 					// diffuse
// 					col := Vec{
// 						rand.Float32() * rand.Float32(),
// 						rand.Float32() * rand.Float32(),
// 						rand.Float32() * rand.Float32(),
// 					}
// 					list[index] = Sphere{center, 0.2, Lambertian{col}}
// 					index++
// 				} else if chooseMat < 0.95 {
// 					// metal
// 					col := Vec{
// 						0.5 * (1 + rand.Float32()),
// 						0.5 * (1 + rand.Float32()),
// 						0.5 * (1 + rand.Float32()),
// 					}
// 					list[index] = Sphere{center, 0.2, Metal{col, 0.3 * rand.Float32()}}
// 					index++
// 				} else {
// 					// glass
// 					list[index] = Sphere{center, 0.2, Dielec{1.5}}
// 					index++
// 				}
// 			}
// 		}
// 	}
// 	list[index] = Sphere{Vec{0, 1, 0}, 1, Dielec{1.5}}
// 	index++
// 	list[index] = Sphere{Vec{-4, 1, 0}, 1, Lambertian{Vec{0.2, 0.5, 0.3}}}
// 	index++
// 	list[index] = Sphere{Vec{4, 1, 0}, 1, Metal{Vec{0.7, 0.6, 0.5}, 0.0}}
// 	index++
// 	return list[:index]
// }

// A central metal sphere surrounded by 6 lambertian, dielectric, and metal spheres
func randomScene() HitList {
	num := 8
	list := make(HitList, num)
	// ground
	list[0] = Sphere{Vec{0, -1000, 0}, 1000, Lambertian{Vec{0.5, 0.5, 0.5}}}
	// centerpiece
	list[1] = Sphere{Vec{0, 1, 0}, 1, Metal{Vec{rand.Float32()*0.5 + 0.5, rand.Float32()*0.5 + 0.5, rand.Float32()*0.5 + 0.5}, 0.0}}
	// Sidepieces
	angle := (math.Pi * 2) / (float64(num) - 2)
	for i := 2; i < num; i++ {
		size := rand.Float32()*0.6 + 0.25
		x := float32(math.Cos(float64(i)*angle)) * 2
		z := float32(math.Sin(float64(i)*angle)) * 2
		col := Vec{rand.Float32()*0.75 + 0.25, rand.Float32()*0.75 + 0.25, rand.Float32()*0.75 + 0.25}
		if i%3 == 0 {
			list[i] = Sphere{Vec{x, size, z}, size, Lambertian{col}}
		} else if i%3 == 1 {
			list[i] = Sphere{Vec{x, size, z}, size, Dielec{rand.Float32()*0.75 + 0.25}}
		} else {
			list[i] = Sphere{Vec{x, size, z}, size, Metal{col, rand.Float32()*0.5 + 0.5}}
		}
	}
	return list[:num]
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
	right := Vec{0.3, 0.5, 1.0}.Mul(t)
	return left.Add(right)
}

// Trace returns a pointer to the traced image
func Trace(args Args, i int, j int, output chan pixel) {
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
	output <- pixel{r: ir, g: ig, b: ib, x: i, y: j}
}

func worker(args Args, input chan pixel, output chan pixel, id int) {
	for i := 0; i < args.Nx*args.Ny; i++ {
		select {
		case work := <-input:
			Trace(args, work.x, work.y, output)
		default:
			return
		}
	}
}

func imgToBase64(img *image.NRGBA) string {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// GoTrace accumulates multiple traces
func GoTrace(seed int, threads int) *image.NRGBA {
	rand.Seed(int64(seed))
	nx := 500
	ny := 500
	ns := 50
	lookFrom := Vec{2.5, 2, 2}
	lookAt := Vec{0, 1.0, 0}
	distToFocus := lookFrom.Sub(lookAt).Len()
	cam := CreateCam(80, float32(nx)/float32(ny), lookFrom, lookAt, Vec{0, 1, 0}, 0.05, distToFocus)
	args := Args{
		World:  randomScene(),
		Nx:     nx,
		Ny:     ny,
		Ns:     ns,
		Camera: cam,
	}
	firstArgs := args
	firstArgs.Ns = 1
	work := make(chan pixel, args.Nx*args.Ny)
	acc := make(chan pixel, args.Nx*args.Ny)
	img := image.NewNRGBA(image.Rect(0, 0, args.Nx, args.Ny))
	for j := 0; j < args.Ny; j++ {
		for i := 0; i < args.Nx; i++ {
			work <- pixel{x: i, y: j}
		}
	}
	for i := 0; i < threads; i++ {
		go worker(args, work, acc, i)
	}
	periodicWrite := time.After(time.Second)
	for c := 0; c < args.Nx*args.Ny; c++ {
		select {
		case p := <-acc:
			img.Set(p.x, args.Ny-p.y-1, color.NRGBA{p.r, p.g, p.b, 255})
		case <-periodicWrite:
			fmt.Println("% done:", float32(c)/float32(args.Nx*args.Ny)*100)
			// send this over websocket or whatever
			// baseString := imgToBase64(img)
			f, err := os.Create("client/img.jpg")
			if err != nil {
				panic(err)
			}
			if err = jpeg.Encode(f, img, nil); err != nil {
				f.Close()
				panic(err)
			}
			periodicWrite = time.After(time.Second)
		}
	}
	return img
}
