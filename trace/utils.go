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

func Sqrt32(f float32) float32 { return float32(math.Sqrt(float64(f))) }
