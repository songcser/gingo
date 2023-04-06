package utils

import "time"

func DateOfString(date string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02", date, loc)
}

func TimeOfString(t string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation(timeFormat, t, loc)
}

func StringToTime(t string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation(time.RFC3339, t, loc)
}

func TimeFormat(t time.Time) string {
	return t.Format(timeFormat)
}

func JsonTimeFormat(t JsonTime) string {
	return time.Time(t).Format(timeFormat)
}
