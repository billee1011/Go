package gutils

import (
	"math/rand"
	"time"
)

// Probability , prob 几率返回true, 1-prob 几率返回 false
func Probability(prob float32) bool {
	if prob <= 0 {
		return false
	}
	if prob >= 1 {
		return true
	}
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Float32() <= prob
}
