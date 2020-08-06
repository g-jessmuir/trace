package trace

import (
	"math"
	"math/rand"
)

func RandUnitSphere() Vec {
	var p Vec
	for {
		p = Vec{rand.Float32(), rand.Float32(), rand.Float32()}
		p.Sub(Vec{1, 1, 1}).Mul(2)
		if p.SqrLen() < 1.0 {
			return p
		}
	}
}

func RandUnitDisk() Vec {
	var p Vec
	for {
		p = Vec{rand.Float32(), rand.Float32(), 0}.Mul(2).Sub(Vec{1, 1, 0})
		if p.Dot(p) < 1 {
			return p
		}
	}
}

func Sqrt32(f float32) float32 { return float32(math.Sqrt(float64(f))) }

func Schlick(cosine float32, refIdx float32) float32 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 *= r0
	return r0 + (1-r0)*float32(math.Pow(float64(1-cosine), 5))
}
