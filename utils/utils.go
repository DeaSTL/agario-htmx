package utils

import (
	"fmt"
	"math/rand"
)

func GenID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz_"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	for i := 0; i < length; i++ {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

func HSVToRGB(h, s, v float64) (r, g, b uint8) {
	var i int
	var f, p, q, t float64

	if s == 0 {
		// Achromatic (grey)
		r = uint8(v * 255)
		g = uint8(v * 255)
		b = uint8(v * 255)
		return
	}

	h /= 60 // sector 0 to 5
	i = int(h)
	f = h - float64(i) // factorial part of h
	p = v * (1 - s)
	q = v * (1 - s*f)
	t = v * (1 - s*(1-f))

	switch i {
	case 0:
		r = uint8(v * 255)
		g = uint8(t * 255)
		b = uint8(p * 255)
	case 1:
		r = uint8(q * 255)
		g = uint8(v * 255)
		b = uint8(p * 255)
	case 2:
		r = uint8(p * 255)
		g = uint8(v * 255)
		b = uint8(t * 255)
	case 3:
		r = uint8(p * 255)
		g = uint8(q * 255)
		b = uint8(v * 255)
	case 4:
		r = uint8(t * 255)
		g = uint8(p * 255)
		b = uint8(v * 255)
	default: // case 5:
		r = uint8(v * 255)
		g = uint8(p * 255)
		b = uint8(q * 255)
	}
	return
}

func GenerateRandomHexColor() string {
	hue := rand.Float64() * 360 // Random hue between 0 and 360
	saturation := 0.8           // Max saturation
	brightness := 1.0           // Max brightness

	r, g, b := HSVToRGB(hue, saturation, brightness)
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}
