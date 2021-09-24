package api

import "time"

func FormatTimestamp(ts time.Time) time.Time {
	loc, _ := time.LoadLocation("Europe/Berlin")
	localTs := time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, loc)
	return localTs
}
