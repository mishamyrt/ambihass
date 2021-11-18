package color

import (
	"image"
	"math"
	"sort"

	"github.com/mishamyrt/ambihass/internal/hass"
)

const downSizeFactor = 250.
const minBucket = .06

type bucket struct {
	Red   float64
	Green float64
	Blue  float64
	Count float64
}

type byCount []bucket

func (c byCount) Len() int           { return len(c) }
func (c byCount) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byCount) Less(i, j int) bool { return c[i].Count < c[j].Count }

func ExtractColors(image image.Image) []hass.RGBColor {
	width := image.Bounds().Max.X
	height := image.Bounds().Max.Y

	stepX := int(math.Max(float64(width)/downSizeFactor, 1))
	stepY := int(math.Max(float64(height)/downSizeFactor, 1))

	var buckets [2][2][2]bucket
	totalCount := 0.
	for x := 0; x < width; x += stepX {
		for y := 0; y < height; y += stepY {
			color := image.At(x, y)
			r, g, b, _ := color.RGBA()
			r >>= 8
			g >>= 8
			b >>= 8
			i := r >> 7
			j := g >> 7
			k := b >> 7
			buckets[i][j][k].Red += float64(r)
			buckets[i][j][k].Green += float64(g)
			buckets[i][j][k].Blue += float64(b)
			buckets[i][j][k].Count += 1
			totalCount++
		}
	}
	var bucketsAverages []bucket
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				currentBucket := buckets[i][j][k]
				if currentBucket.Count > 0 {
					bucketsAverages = append(bucketsAverages, bucket{
						Count: currentBucket.Count,
						Red:   currentBucket.Red / currentBucket.Count,
						Green: currentBucket.Green / currentBucket.Count,
						Blue:  currentBucket.Blue / currentBucket.Count,
					})
				}
			}
		}
	}

	sort.Sort(sort.Reverse(byCount(bucketsAverages)))
	colors := []hass.RGBColor{}
	for _, bucket := range bucketsAverages {
		if bucket.Count/totalCount > minBucket {
			colors = append(colors, hass.RGBColor{
				uint32(uint8(math.Round(bucket.Red))),
				uint32(uint8(math.Round(bucket.Green))),
				uint32(uint8(math.Round(bucket.Blue))),
			})
		}
	}

	return colors
}
