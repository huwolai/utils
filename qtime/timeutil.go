package qtime

import (
	"time"
)

func ToyyyyMMddHHmm(tm time.Time) string {

	return tm.Format("2006-01-02 15:04")
}


func ToyyyyMM2(tm time.Time) string {

	return tm.Format("200601")
}

func ToyyyyMMdd(tm time.Time) string {

	return tm.Format("20060102")
}

//带-的日期
func Toyyyy_MM_dd(tm time.Time) string {

	return tm.Format("2006-01-02")
}

func PareTimeStrForYYYY_mm_dd(timeStr string) (time.Time,error)  {
	return time.Parse("2006-01-02",timeStr)
}


