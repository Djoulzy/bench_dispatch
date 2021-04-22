package datamodels

import (
	"bench_dispatch/clog"
	"time"
)

const MySQL_DateFormat = "2006-01-02 15:04:05"
const IOS_DateFormat = "2006-01-02T15:04:05.00Z"

var Location *time.Location

// DateConvertISOtoMySQL : Converti une date au format UTC iOS vers le format simple MySQL
func DateConvertIOStoMySQL(iosDate string) string {
	var tmp time.Time
	var err error

	if tmp, err = time.ParseInLocation(IOS_DateFormat, iosDate, Location); err != nil {
		clog.Warn("RideManager", "NewRide", "Bad date format %s", iosDate)
		return ""
	}
	return tmp.Format(MySQL_DateFormat)
}

func GetDateFromIOS(date string) (time.Time, error) {
	var tmp time.Time
	var err error
	if tmp, err = time.ParseInLocation(IOS_DateFormat, date, Location); err != nil {
		clog.Warn("RideManager", "NewRide", "Bad date format %s", date)
		return tmp, err
	}
	return tmp, nil
}

func FormatDateForIOS(date time.Time) string {
	return date.Format(IOS_DateFormat)
}

func FormatDateForMySQL(date time.Time) string {
	return date.Format(MySQL_DateFormat)
}
