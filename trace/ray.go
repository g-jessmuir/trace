package trace

// Ray implements a 3D line consisting of an origin and direction
type Ray struct {
	Origin, Dir Vec
}

// PatT returns the point Origin + t*Dir
func (r Ray) PatT(t float32) Vec {
	return r.Origin.Add(r.Dir.Mul(t))
}
