package worldmap

import (
	"math"

	opensimplex "github.com/ojrac/opensimplex-go"
)

func NewWorldMap(seed string, octaves int, lacunarity float64, gain float64) *WorldMap {
	seedInt := hash(seed)
	worldmap := &WorldMap{}
	worldmap.noise = opensimplex.NewWithSeed(seedInt)
	worldmap.fmb = &fractalBrownianMotionOptions{octaves, lacunarity, gain}
	return worldmap
}

type fractalBrownianMotionOptions struct {
	octaves    int
	lacunarity float64
	gain       float64
}

type Region struct {
	X             int
	Y             int
	Height        float64
	LandscapeType string
}

type WorldMap struct {
	noise *opensimplex.Noise
	fmb   *fractalBrownianMotionOptions
}

func (w *WorldMap) GetRegion(x int, y int) *Region {
	region := &Region{}
	region.X = x
	region.X = y
	region.Height = math.Abs(w.fractalBrownianMotion(float64(x)*0.0005, float64(y)*0.0005)) / 2
	region.LandscapeType = w.computeLandscapeType(region.Height)
	return region
}

func (w *WorldMap) fractalBrownianMotion(x float64, y float64) float64 {
	sum := 0.0
	amplitude := 1.0
	for i := 0; i < w.fmb.octaves; i++ {
		sum += amplitude * w.noise.Eval2(x, y)
		amplitude *= w.fmb.gain
		x *= w.fmb.lacunarity
		y *= w.fmb.lacunarity
	}
	return sum
}

func (w *WorldMap) computeLandscapeType(height float64) string {
	switch {
	case height < 0.3:
		return "deepwater"
	case height > 0.3 && height <= 0.35:
		return "water"
	case height > 0.35 && height <= 0.95:
		return "plain"
	case height > 0.95:
		return "snow"
	}
	return "unknown"
}

func hash(str string) int64 {
	var h int64 = 0
	for _, r := range str {
		h = 31*h + int64(r)
	}
	return h
}
