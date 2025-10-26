package main

type Interval struct {
	min, max float64
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
