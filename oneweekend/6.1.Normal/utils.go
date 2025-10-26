package main

import (
	"math"
)

func Clamp[T int | float64](val T, min T, max T) T {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
