package trace

type Mat interface {
	// TODO: change scatter function to just return attenuation and scattered
	Scatter(rIn Ray, rec *HitRecord, attenuation *Vec, scattered *Ray) bool
}

type Lambertian struct {
	Albedo Vec
}

func (l Lambertian) Scatter(rIn Ray, rec *HitRecord, attenuation *Vec, scattered *Ray) bool {
	target := rec.P.Add(rec.N).Add(RandUnitSphere())
	*scattered = Ray{Origin: rec.P, Dir: target.Sub(rec.P)}
	*attenuation = l.Albedo
	return true
}

type Metal struct {
	Albedo Vec
	Fuzz   float32
}

func (m Metal) Scatter(rIn Ray, rec *HitRecord, attenuation *Vec, scattered *Ray) bool {
	reflected := rIn.Dir.Unit().Reflect(rec.N)
	*scattered = Ray{Origin: rec.P, Dir: reflected.Add(RandUnitSphere().Mul(m.Fuzz))}
	*attenuation = m.Albedo
	return scattered.Dir.Dot(rec.N) > 0
}

type Dielec struct {
	RefIdx float32
}

func (d Dielec) Scatter(rIn Ray, rec *HitRecord, attenuation *Vec, scattered *Ray) bool {
	var outNorm Vec
	reflected := rIn.Dir.Reflect(rec.N)
	var niByNt float32
	*attenuation = Vec{1, 1, 1}
	if rIn.Dir.Dot(rec.N) > 0 {
		outNorm = Vec{0, 0, 0}.Sub(rec.N)
		niByNt = d.RefIdx
	} else {
		outNorm = rec.N
		niByNt = 1 / d.RefIdx
	}
	if refracted, didHit := rIn.Dir.Refract(outNorm, niByNt); didHit {
		*scattered = Ray{rec.P, refracted}
	} else {
		*scattered = Ray{rec.P, reflected}
		return false
	}
	return true
}
