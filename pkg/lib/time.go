package lib

import (
	"strings"
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/utils"
)

// CurrentTime func
// will return in format "2006-01-02 15:04:05"
func CurrentTime(format ...string) string {

	form := TimeFormat()
	if len(format) > 0 {
		form = format[0]
	}
	return time.Now().Format(form)
}

// TimeNow func
func TimeNow() time.Time {
	return time.Now()
}

func TimeWIB(d ...time.Time) time.Time {
	t := TimeNow()
	if len(d) > 0 {
		t = d[0]
	}
	return t.Add(7 * time.Hour)
}

// TimeNow func
func TimeNowPtr() *time.Time {
	return Pointer(time.Now())
}

func TimeFormat() string {
	return "2006-01-02 15:04:05"
	// return time.RFC3339
}

func SetTime(date time.Time, hour, minute, second int) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, second, 0, time.UTC)
}

func ExpireNotif(pushAt time.Time) time.Duration {
	return pushAt.Add(5 * time.Minute).Sub(TimeNow())
}

func IsStarted(t time.Time, d ...time.Time) bool {
	x := TimeNow()
	if len(d) > 0 {
		x = d[0]
	}
	return t.Before(x)
}

func InDate(date, check time.Time) bool {
	return check.After(SetTime(date, 0, 0, 1)) && check.Before(SetTime(date, 23, 59, 59))
}

func AddDate(t time.Time, day int) time.Time {
	return t.AddDate(0, 0, day)
}

func AddMonth(t time.Time, month int) time.Time {
	return t.AddDate(0, month, 0)
}

func AddYear(t time.Time, year int) time.Time {
	return t.AddDate(year, 0, 0)
}

func CompareTime(a, b time.Time) bool {
	if a.IsZero() && b.IsZero() {
		return true
	}
	if a.IsZero() || b.IsZero() {
		return false
	}
	utils.Debug(a.Format(TimeFormat()), b.Format(TimeFormat()))
	return strings.EqualFold(a.Format(TimeFormat()), b.Format(TimeFormat()))
}
