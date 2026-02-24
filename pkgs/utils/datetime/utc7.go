package datetime

import "time"

var _local *time.Location

func GetTimeNowWithUTC7() time.Time {
	if _local == nil {
		loc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			panic("Failed to load location: " + err.Error())
		}
		_local = loc
	}
	return time.Now().In(_local)
}
