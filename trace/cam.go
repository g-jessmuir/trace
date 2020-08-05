package trace

type Cam struct {
	Origin     Vec
	LowerLeft  Vec
	Horizontal Vec
	Vertical   Vec
}

func (c Cam) GetRay(u float32, v float32) Ray {
	return Ray{
		Origin: c.Origin,
		Dir:    c.LowerLeft.Add(c.Horizontal.Mul(u)).Add(c.Vertical.Mul(v)).Sub(c.Origin),
	}
}
