package perlin

import (
	"math"
	"math/rand"
    "image"
    "image/color"
    "image/png"
    "os"
)

type Generator struct {
    permutations []int
}

func NewGenerator(seed int64) Generator {
    r := rand.New(rand.NewSource(seed))
    return Generator { r.Perm(256) }
}

func (g Generator) Noise1D(x float64) float64 {
    return g.Noise2D(x, 0.5)
}

func (g Generator) Noise2D(x, y float64) float64 {
	return g.Noise3D(x, y, 0.5)
}

func (g Generator) Noise3D(x, y, z float64) float64 {
	X := int(math.Floor(x)) & 255
	Y := int(math.Floor(y)) & 255
	Z := int(math.Floor(z)) & 255

    x, y, z = x-math.Floor(x), y-math.Floor(y), z-math.Floor(z)
	u, v, w := fade(x), fade(y), fade(z)

	A, B := g.p(X)+Y, g.p(X+1)+Y
	AA, AB, BA, BB := g.p(A)+Z, g.p(A+1)+Z, g.p(B)+Z, g.p(B+1)+Z

	return (lerp(
		w,
		lerp(
			v,
			lerp(u, grad(g.p(AA), x, y, z), grad(g.p(BA), x-1, y, z)),
			lerp(u, grad(g.p(AB), x, y-1, z), grad(g.p(BB), x-1, y-1, z)),
		),
		lerp(
			v,
			lerp(u, grad(g.p(AA+1), x, y, z-1), grad(g.p(BA+1), x-1, y, z-1)),
			lerp(u, grad(g.p(AB+1), x, y-1, z-1), grad(g.p(BB+1), x-1, y-1, z-1)),
		),
	))
}

func (g Generator) p(idx int) int {
    return g.permutations[idx % 255]
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y, z float64) float64 {
	var h int = hash & 15

	var u float64
	if h < 8 {
		u = x
	} else {
		u = y
	}

	var v float64
	if h < 4 {
		v = y
	} else {
		if h == 12 || h == 14 {
			v = x
		} else {
			v = z
		}
	}

	var r1 float64
	if (h & 1) == 0 {
		r1 = u
	} else {
		r1 = -u
	}

	var r2 float64
	if (h & 2) == 0 {
		r2 = v
	} else {
		r2 = -v
	}

	return r1 + r2
}
