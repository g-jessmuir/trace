package trace

import "math"

type Cam struct {
	Origin     Vec
	LowerLeft  Vec
	Horizontal Vec
	Vertical   Vec
	u, v, w    Vec
	LensRadius float32
}

func (c Cam) getRay(s float32, t float32) Ray {
	rd := RandUnitDisk().Mul(c.LensRadius)
	offset := c.u.Mul(rd.X).Add(c.v.Mul(rd.Y))
	origin := c.Origin.Add(offset)
	dir := c.LowerLeft.Add(c.Horizontal.Mul(s)).Add(c.Vertical.Mul(t)).Sub(c.Origin).Sub(offset)
	return Ray{
		Origin: origin,
		Dir:    dir,
	}
}

func CreateCam(vfov float32, aspect float32, from Vec, at Vec, up Vec, aperture float32,
	focusDist float32) Cam {
	theta := vfov * math.Pi / 180
	halfHeight := float32(math.Tan(float64(theta / 2)))
	halfWidth := aspect * halfHeight
	w := from.Sub(at).Unit()
	u := up.Cross(w).Unit()
	v := w.Cross(u)
	return Cam{
		Origin:     from,
		LowerLeft:  from.Sub(u.Mul(halfWidth * focusDist)).Sub(v.Mul(halfHeight * focusDist)).Sub(w.Mul(focusDist)),
		Horizontal: u.Mul(2 * halfWidth * focusDist),
		Vertical:   v.Mul(2 * halfHeight * focusDist),
		u:          u,
		v:          v,
		w:          w,
		LensRadius: aperture / 2,
	}
}
