package trace

// Ray implements a 3D line consisting of an origin and direction
type Ray struct {
	A, B Vec
}

// Origin returns the ray's origin point
func (r Ray) Origin() Vec { return r.A }

// Dir returns the ray's direction vector
func (r Ray) Dir() Vec { return r.B }

// PatT returns the point A + t*B
func (r Ray) PatT(t float32) Vec {
	return r.A.Add(r.B.Mul(t))
}
