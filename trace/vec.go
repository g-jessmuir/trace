package trace

import "math"

// Vec implements a Vector type for vector math
type Vec struct {
	X, Y, Z float32
}

// Add o to v and return a new vector with the result.
func (v Vec) Add(o Vec) Vec {
	return Vec{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

// Sub o from v and return a new vector with the result.
func (v Vec) Sub(o Vec) Vec {
	return Vec{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

// Mul multiplies v by f and returns the result
func (v Vec) Mul(f float32) Vec {
	return Vec{v.X * f, v.Y * f, v.Z * f}
}

// Div divides v by f and returns the result
func (v Vec) Div(f float32) Vec {
	return Vec{v.X / f, v.Y / f, v.Z / f}
}

// VDiv divides v by o component-wise and returns the result
func (v Vec) VDiv(o Vec) Vec {
	return Vec{v.X / o.X, v.Y / o.Y, v.Z / o.Z}
}

// Adda adds o to v and assigns the result to v
func (v Vec) Adda(o *Vec) {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
}

// Suba subtracts o from v and assigns the result to v
func (v *Vec) Suba(o *Vec) {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
}

// Mula mutliples v by o component-wise and assigns the result to v
func (v *Vec) Mula(o *Vec) {
	v.X *= o.X
	v.Y *= o.Y
	v.Z *= o.Z
}

// Diva divides v by o component-wise and assigns the result to v
func (v *Vec) Diva(o *Vec) {
	v.X /= o.X
	v.Y /= o.Y
	v.Z /= o.Z
}

// SMula multiplies v by scalar f and assigns the result to v
func (v *Vec) SMula(f float32) {
	v.X *= f
	v.Y *= f
	v.Z *= f
}

// SDiva divides v by scalar f and assigns the result to v
func (v *Vec) SDiva(f float32) {
	v.X /= f
	v.Y /= f
	v.Z /= f
}

// Len returns the length of v
func (v Vec) Len() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// SqrLen returns the squared length of v
func (v Vec) SqrLen() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Unit returns a new unit vector of v
func (v Vec) Unit() Vec {
	l := v.Len()
	v.SDiva(l)
	return v
}

func (v Vec) Dot(o Vec) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}
