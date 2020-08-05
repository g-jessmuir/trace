package trace

import "math"

type HitRecord struct {
	T float32
	P Vec
	N Vec
}

type Hitable interface {
	Hit(r Ray, tMin float32, tMax float32, rec *HitRecord) bool
}

type HitList []Hitable

type Sphere struct {
	Center Vec
	Radius float32
}

func (hs HitList) Hit(r Ray, tMin float32, tMax float32, rec *HitRecord) bool {
	var hr HitRecord
	var hitAnything bool
	closestSoFar := tMax
	for _, h := range hs {
		if h.Hit(r, tMin, closestSoFar, &hr) {
			hitAnything = true
			closestSoFar = hr.T
			*rec = hr
		}
	}
	return hitAnything
}

func (s Sphere) Hit(r Ray, tMin float32, tMax float32, rec *HitRecord) bool {
	oc := r.Origin().Sub(s.Center)
	a := r.Dir().Dot(r.Dir())
	b := oc.Dot(r.Dir())
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c
	if discriminant > 0 {
		p := (-b - float32(math.Sqrt(float64(discriminant)))) / a
		if p < tMax && p > tMin {
			rec.T = p
			rec.P = r.PatT(rec.T)
			rec.N = rec.P.Sub(s.Center).Div(s.Radius)
			return true
		}
		p = (-b + float32(math.Sqrt(float64(discriminant)))) / a
		if p < tMax && p > tMin {
			rec.T = p
			rec.P = r.PatT(rec.T)
			rec.N = rec.P.Sub(s.Center).Div(s.Radius)
			return true
		}
	}
	return false
}
