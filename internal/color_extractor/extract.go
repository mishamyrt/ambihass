package color_extractor

import (
	"image"
	"math"
	"sort"

	"github.com/mishamyrt/ambihass/internal/hass"
)

type bucket struct {
	Red   float64
	Green float64
	Blue  float64
	Count float64
}

type ByCount []bucket

func (c ByCount) Len() int           { return len(c) }
func (c ByCount) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCount) Less(i, j int) bool { return c[i].Count < c[j].Count }

type Config struct {
	DownSizeTo  float64
	SmallBucket float64
}

func ExtractColors(image image.Image) []hass.RGBColor {
	return ExtractColorsWithConfig(image, Config{
		DownSizeTo:  300.,
		SmallBucket: .01,
	})
}

func ExtractColorsWithConfig(image image.Image, config Config) []hass.RGBColor {
	width := image.Bounds().Max.X
	height := image.Bounds().Max.Y

	stepX := int(math.Max(float64(width)/config.DownSizeTo, 1))
	stepY := int(math.Max(float64(height)/config.DownSizeTo, 1))

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

	sort.Sort(sort.Reverse(ByCount(bucketsAverages)))

	colors := []hass.RGBColor{}
	for _, avg := range bucketsAverages {
		if avg.Count/totalCount > config.SmallBucket {
			colors = append(colors, hass.RGBColor{
				uint32(uint8(math.Round(avg.Red))),
				uint32(uint8(math.Round(avg.Green))),
				uint32(uint8(math.Round(avg.Blue))),
			})
		}
	}

	return colors
}
