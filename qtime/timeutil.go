package qtime

import (
	"time"
)

func ToyyyyMMddHHmm(tm time.Time) string {

	return tm.Format("2006-01-02 15:04")
}



