package main

import "math"

type Interval struct {
	min, max float64
}

var EmptyInterval = Interval{math.MaxFloat64, -math.MaxFloat64}
var UniverseInterval = Interval{-math.MaxFloat64, math.MaxFloat64}

func NewInterval(a, b Interval) Interval {
	min := b.min
	if a.min <= b.min {
		min = a.min
	}
	max := b.max
	if a.max >= b.max {
		max = a.max
	}
	return Interval{min, max}
}

func (i Interval) Size() float64 {
	return i.max - i.min
}

func (i Interval) Contains(x float64) bool {
	return x >= i.min && x <= i.max
}

func (i Interval) Surrounds(x float64) bool {
	return x > i.min && x < i.max
}

func (i Interval) Clamp(x float64) float64 {
	return Clamp(x, i.min, i.max)
}

func (i Interval) Expand(delta float64) Interval {
	padding := delta / 2
	return Interval{i.min - padding, i.max + padding}
}
