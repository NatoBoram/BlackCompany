package main

import "github.com/aiseeq/s2l/lib/point"

func Midpoint(p1, p2 point.Point) point.Point {
	return (p1 + p2).Mul(0.5)
}
