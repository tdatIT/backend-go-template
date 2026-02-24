package datetime

import (
	"errors"
	"time"
)

func ParseDatetimeFromString(timeStr string) (time.Time, error) {
	value, err := time.Parse(time.DateTime, timeStr)
	return value, err
}

func ParseTimeWithUTC7(datetimeStr string) (time.Time, error) {
	const layout = "2006-01-02 15:04:05"

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Time{}, errors.New("cannot load Asia/Bangkok timezone")
	}

	parsedTime, err := time.ParseInLocation(layout, datetimeStr, location)
	if err != nil {
		return time.Time{}, errors.New("invalid datetime format, expected yyyy-mm-dd hh:mm:ss")
	}

	return parsedTime, nil
}

func ParseDateFromString(dateStr string) (time.Time, error) {
	return time.Parse(time.DateOnly, dateStr)
}

func ParseTimeFromString(timeStr string) (time.Time, error) {
	return time.Parse(time.TimeOnly, timeStr)
}

func FormatDatetime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
func ParseDateTimeToDateString(t time.Time) string {
	return t.Format(time.DateOnly)
}
func FormatDatetimeWithUtc7(t time.Time) string {
	return t.In(time.FixedZone("UTC+7", 7*3600)).Format("2006-01-02 15:04:05")
}

func FormatDate(time time.Time) string {
	return time.Format("2006-01-02")
}

func FormatTime(time time.Time) string {
	return time.Format("15:04:05")
}

func FormatDateTimeVN(time time.Time) string {
	return time.Format("02/01/2006 15:04:05")
}

func ConvertSQLDateFormatToDate(dateStr string) string {
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ""
	}
	return parsedTime.Format("02/01/2006")
}

func ConvertUTCToAsiaHCM(utcTime time.Time) time.Time {
	return utcTime.In(time.FixedZone("UTC+7", 7*3600))
}

func ParseDateTimeWithUTC7(datetimeStr string) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return time.Time{}, errors.New("cannot load Asia/Bangkok timezone")
	}

	parsedTime, err := time.ParseInLocation(time.DateTime, datetimeStr, location)
	if err != nil {
		return time.Time{}, errors.New("invalid datetime format, expected yyyy-mm-dd hh:mm:ss")
	}

	return parsedTime, nil
}
