// timegetter
package service

import (
	"time"
)

func GetTime() int64 {
	timestamp := time.Now().Unix()
	return timestamp
}
