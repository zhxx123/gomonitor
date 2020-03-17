package utils

import (
	"math/rand"
	"time"
)

func GetRandNumInt(maxnum int) int {
	rand.Seed(time.Now().UnixNano())
	ikind := rand.Intn(maxnum)
	return ikind
}
